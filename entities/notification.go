package entities

import "time"

// Notification represents a persisted in-app notification.
// The `type` field drives deep-link routing in Flutter:
//   - "assignment"   → tugas page
//   - "presence"     → presensi page
//   - "announcement" → general info
//   - "general"      → no redirect
type Notification struct {
	ID            int64      `gorm:"column:id;primaryKey"             json:"id"`
	RecipientID   string     `gorm:"column:recipient_id"              json:"recipient_id"`
	RecipientRole string     `gorm:"column:recipient_role"            json:"recipient_role"`
	SenderID      string     `gorm:"column:sender_id"                 json:"sender_id"`
	SenderName    string     `gorm:"column:sender_name"               json:"sender_name"`
	Title         string     `gorm:"column:title"                     json:"title"`
	Body          string     `gorm:"column:body"                      json:"body"`
	Type          string     `gorm:"column:type"                      json:"type"`
	ReferenceID   string     `gorm:"column:reference_id"              json:"reference_id"`
	IsRead        bool       `gorm:"column:is_read;default:false"     json:"is_read"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	ReadAt        *time.Time `gorm:"column:read_at"                   json:"read_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
