package entities

type Major struct {
	MajorID          string `gorm:"column:jurid" json:"major_id"`
	RealMajorID      string `gorm:"column:jurasli_jurid" json:"real_major_id"`
	MajorName        string `gorm:"column:jurnama" json:"major_name"`
	StudyProgramName string `gorm:"column:prodinama" json:"study_program_name"`
	// FacultyID        string `gorm:"column:fakid"`
	// Title                    string `gorm:"column:gelar"`
	// FullTitle                string `gorm:"column:gelarpjng"`
	// AccreditationCertificate string `gorm:"column:sk_akreditasi"`
	// StudyProgramID           string `gorm:"column:id_program_studi"`
	// WorkUnitID               int    `gorm:"column:unitkerjaid"`
	// MajorParentID            string `gorm:"column:jur_parent_id"`
	// IsActive                 bool   `gorm:"column:isactive"`
}
