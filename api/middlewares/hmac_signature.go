package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// HmacSignature memverifikasi setiap request dari mobile app.
//
// Header yang wajib ada:
//   X-Timestamp : Unix timestamp detik (string), mis. "1714285351"
//   X-Signature : HMAC-SHA256 hex dari payload string berikut:
//                 "<METHOD>\n<PATH>\n<TIMESTAMP>\n<BODY>"
//
// Contoh payload yang di-sign:
//   POST\n/api/v1/students/subjects\n1714285351\n{"academic_period":"..."}
//
// Toleransi waktu: ±5 menit (anti replay-attack)
func HmacSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsHeader := c.GetHeader("X-Timestamp")
		sigHeader := c.GetHeader("X-Signature")

		// 1. Header wajib ada
		if tsHeader == "" || sigHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Missing security headers (X-Timestamp, X-Signature)",
			})
			c.Abort()
			return
		}

		// 2. Parse timestamp
		ts, err := strconv.ParseInt(tsHeader, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "X-Timestamp tidak valid",
			})
			c.Abort()
			return
		}

		// 3. Cek toleransi waktu ±5 menit (anti replay-attack)
		diff := math.Abs(float64(time.Now().Unix() - ts))
		if diff > 300 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Request kadaluarsa (timestamp terlalu jauh dari waktu server)",
			})
			c.Abort()
			return
		}

		// 4. Baca body (gin memungkinkan re-read body)
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Gagal membaca request body",
			})
			c.Abort()
			return
		}
		// Restore body agar handler berikutnya bisa baca ulang
		c.Request.Body = io.NopCloser(bytesReader(bodyBytes))

		// 5. Susun payload yang sama persis dengan Flutter
		payload := fmt.Sprintf("%s\n%s\n%s\n%s",
			c.Request.Method,
			c.Request.URL.Path,
			tsHeader,
			string(bodyBytes),
		)

		// 6. Hitung HMAC-SHA256
		secret := viper.GetString("SECRET_KEY")
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(payload))
		expected := hex.EncodeToString(mac.Sum(nil))

		// 7. Bandingkan secara constant-time (mencegah timing attack)
		//
		// ⚠️  MODE SAAT INI: warn-only
		// Ubah enforceHmac = true setelah semua repository Flutter
		// sudah menggunakan ApiClient dengan HmacInterceptor.
		const enforceHmac = true

		if !hmac.Equal([]byte(sigHeader), []byte(expected)) {
			if enforceHmac {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "failed",
					"message": "Signature tidak valid",
				})
				c.Abort()
				return
			}
			// warn-only: lanjut tapi log warning
			log.Printf("[HMAC] ⚠️  signature tidak cocok — path=%s method=%s ts=%s",
				c.Request.URL.Path, c.Request.Method, tsHeader)
		}

		c.Next()
	}
}

// bytesReader adalah helper agar body bisa dibaca ulang.
type bytesReaderWrapper struct {
	data []byte
	pos  int
}

func bytesReader(b []byte) io.Reader {
	return &bytesReaderWrapper{data: b}
}

func (r *bytesReaderWrapper) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
