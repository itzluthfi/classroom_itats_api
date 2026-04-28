package microsoft_services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"classroom_itats_api/entities"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// ─── Response Structs ────────────────────────────────────────────────────────

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

type OnlineMeeting struct {
	ID          string `json:"id"`
	JoinWebURL  string `json:"joinWebUrl"`
	Subject     string `json:"subject"`
	StartTime   string `json:"startDateTime"`
	EndTime     string `json:"endDateTime"`
}

// ─── Service ─────────────────────────────────────────────────────────────────

type MicrosoftAuthService struct {
	db *gorm.DB
}

func NewMicrosoftAuthService(db *gorm.DB) *MicrosoftAuthService {
	// Auto-migrate tabel ms_tokens jika belum ada
	db.AutoMigrate(&entities.MicrosoftToken{})
	return &MicrosoftAuthService{db: db}
}

// GetOAuthURL menghasilkan URL untuk redirect dosen ke halaman login Microsoft
func (s *MicrosoftAuthService) GetOAuthURL(state string) string {
	tenantID := viper.GetString("AZURE_TENANT_ID")
	clientID := viper.GetString("AZURE_CLIENT_ID")
	redirectURI := viper.GetString("AZURE_REDIRECT_URI")

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("response_mode", "query")
	params.Set("scope", "offline_access OnlineMeetings.ReadWrite")
	params.Set("state", state)

	return fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?%s",
		tenantID,
		params.Encode(),
	)
}

// ExchangeCodeForToken menukar auth code menjadi access+refresh token
func (s *MicrosoftAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	tenantID := viper.GetString("AZURE_TENANT_ID")
	clientID := viper.GetString("AZURE_CLIENT_ID")
	clientSecret := viper.GetString("AZURE_CLIENT_SECRET")
	redirectURI := viper.GetString("AZURE_REDIRECT_URI")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", "offline_access OnlineMeetings.ReadWrite")

	return s.postToTokenEndpoint(ctx, tenantID, data)
}

// RefreshAccessToken memperbarui access token menggunakan refresh token yang tersimpan
func (s *MicrosoftAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	tenantID := viper.GetString("AZURE_TENANT_ID")
	clientID := viper.GetString("AZURE_CLIENT_ID")
	clientSecret := viper.GetString("AZURE_CLIENT_SECRET")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("scope", "offline_access OnlineMeetings.ReadWrite")

	return s.postToTokenEndpoint(ctx, tenantID, data)
}

// SaveToken menyimpan atau memperbarui token di database berdasarkan dosID
func (s *MicrosoftAuthService) SaveToken(ctx context.Context, dosID string, token *TokenResponse) error {
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second).Unix()

	result := s.db.WithContext(ctx).Where(entities.MicrosoftToken{DosID: dosID}).
		Assign(entities.MicrosoftToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    expiresAt,
		}).
		FirstOrCreate(&entities.MicrosoftToken{})

	return result.Error
}

// GetValidAccessToken mengambil access token yang valid (refresh otomatis jika expired)
func (s *MicrosoftAuthService) GetValidAccessToken(ctx context.Context, dosID string) (string, error) {
	var msToken entities.MicrosoftToken
	if err := s.db.WithContext(ctx).Where("dos_id = ?", dosID).First(&msToken).Error; err != nil {
		return "", fmt.Errorf("dosen belum pernah login Microsoft, silakan hubungkan akun terlebih dahulu")
	}

	// Cek apakah token masih valid (buffer 5 menit)
	if time.Now().Unix() < msToken.ExpiresAt-300 {
		return msToken.AccessToken, nil
	}

	// Token expired → refresh
	newToken, err := s.RefreshAccessToken(ctx, msToken.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("sesi Microsoft telah berakhir, silakan login ulang: %w", err)
	}

	// Simpan token baru
	if err := s.SaveToken(ctx, dosID, newToken); err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

// CreateOnlineMeeting membuat online meeting di MS Teams via Graph API
func (s *MicrosoftAuthService) CreateOnlineMeeting(ctx context.Context, accessToken, subject, startTime, endTime string) (*OnlineMeeting, error) {
	payload := map[string]interface{}{
		"subject": subject,
		"startDateTime": startTime, // format: "2024-01-01T09:00:00Z"
		"endDateTime":   endTime,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://graph.microsoft.com/v1.0/me/onlineMeetings",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MS Graph API error [%d]: %s", resp.StatusCode, string(respBody))
	}

	var meeting OnlineMeeting
	if err := json.Unmarshal(respBody, &meeting); err != nil {
		return nil, err
	}

	return &meeting, nil
}

// HasLinkedAccount mengecek apakah dosen sudah pernah menghubungkan akun Microsoft
func (s *MicrosoftAuthService) HasLinkedAccount(ctx context.Context, dosID string) bool {
	var count int64
	s.db.WithContext(ctx).Model(&entities.MicrosoftToken{}).Where("dos_id = ?", dosID).Count(&count)
	return count > 0
}

// ─── Private helpers ──────────────────────────────────────────────────────────

func (s *MicrosoftAuthService) postToTokenEndpoint(ctx context.Context, tenantID string, data url.Values) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed [%d]: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
