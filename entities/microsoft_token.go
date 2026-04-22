package entities

import "time"

// MicrosoftToken menyimpan OAuth token MS per dosen
// Tabel: ms_tokens (di-migrate via GORM AutoMigrate)
type MicrosoftToken struct {
	ID           uint      `gorm:"primarykey;autoIncrement" json:"id"`
	DosID        string    `gorm:"column:dos_id;uniqueIndex;not null" json:"dos_id"`
	AccessToken  string    `gorm:"column:access_token;type:text;not null" json:"-"`
	RefreshToken string    `gorm:"column:refresh_token;type:text;not null" json:"-"`
	ExpiresAt    int64     `gorm:"column:expires_at;not null" json:"expires_at"` // Unix timestamp
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MicrosoftToken) TableName() string {
	return "ms_tokens"
}
