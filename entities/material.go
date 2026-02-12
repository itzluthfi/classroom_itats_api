package entities

type Material struct {
	MaterialID        string `gorm:"column:materi_id" json:"material_id"`
	LecturerID        string `gorm:"column:dosid" json:"lecturer_id"`
	MaterialTitle     string `gorm:"column:judul_materi" json:"material_title"`
	MaterialLink      string `gorm:"column:link_materi" json:"material_link"`
	CreatedAt         string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         string `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         string `gorm:"column:deleted_at" json:"deleted_at"`
	HiddenStatus      int16  `gorm:"column:hidden_status" json:"hidden_status"`
	LectureMaterialID string `gorm:"column:kul_materi_id" json:"lecture_material_id"`
	LectureID         string `gorm:"column:kul_id" json:"lecture_id"`
}
