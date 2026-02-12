package entities

import "time"

type AnnouncementJSON struct {
	AnnouncementID int                    `json:"announcement_id"`
	PostContent    string                 `json:"post_content"`
	CreatedAt      time.Time              `json:"created_at"`
	AuthorID       string                 `json:"author_id"`
	Author         string                 `json:"author"`
	Photo          string                 `json:"photo"`
	Materials      []AnnouncementMaterial `json:"materials"`
	Comments       []Comment              `json:"comments"`
}

type Announcement struct {
	AnnouncementID int       `gorm:"type:integer;column:id_post_klstw"`
	PostContent    string    `gorm:"type:jsonb;column:content_post"`
	CreatedAt      time.Time `gorm:"type:timestamp;column:created_at"`
	AuthorID       string    `gorm:"type:character varying(50);column:author_id"`
	Author         string    `gorm:"type:character varying(50);column:nama"`
	Photo          string    `gorm:"type:character varying(100);column:foto"`
	Comments       string    `gorm:"type:jsonb;column:comments"`
	Materials      string    `gorm:"type:jsonb;column:post_materi"`
}

type AnnouncementMaterial struct {
	MaterialID        string `json:"materi_id"`
	LecturerID        string `json:"dosid"`
	MaterialTitle     string `json:"judul_materi"`
	MaterialLink      string `json:"link_materi"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	DeletedAt         string `json:"deleted_at"`
	HiddenStatus      int16  `json:"hidden_status"`
	LectureMaterialID string `json:"kul_materi_id"`
	LectureID         string `json:"kul_id"`
}

type Comment struct {
	CommentID      int    `json:"id_post_comment"`
	AnnouncementID int    `json:"post_klstw_id"`
	CommentContent string `json:"content_comment"`
	CreatedAt      string `json:"created_at"`
	AuthorID       string `json:"author_id"`
	Author         string `json:"nama"`
	Photo          string `json:"foto"`
}

type AnnouncementStore struct {
	ActivityMasterId string    `gorm:"column:master_kegiatan_id" json:"activity_master_id"`
	PostContent      string    `gorm:"type:jsonb;column:content_post" json:"post_content"`
	AuthorId         string    `gorm:"column:author_id"`
	FlagAuthor       int       `gorm:"column:flag_author"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type AnnouncementUpdate struct {
	AnnouncementID   int       `gorm:"column:id_post_klstw" json:"announcement_id"`
	ActivityMasterId string    `gorm:"column:master_kegiatan_id" json:"activity_master_id"`
	PostContent      string    `gorm:"type:jsonb;column:content_post" json:"post_content"`
	AuthorId         string    `gorm:"column:author_id"`
	FlagAuthor       int       `gorm:"column:flag_author"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CommentStore struct {
	AnnouncementID int       `gorm:"column:post_klstw_id" json:"announcement_id"`
	CommentContent string    `gorm:"column:content_comment" json:"comment_content"`
	AuthorId       string    `gorm:"column:author_id"`
	FlagAuthor     int       `gorm:"column:flag_author"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CommentUpdate struct {
	CommentID      int       `gorm:"column:id_post_comment" json:"comment_id"`
	AnnouncementID int       `gorm:"column:post_klstw_id" json:"announcement_id"`
	CommentContent string    `gorm:"column:content_comment" json:"comment_content"`
	AuthorId       string    `gorm:"column:author_id"`
	FlagAuthor     int       `gorm:"column:flag_author"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}
