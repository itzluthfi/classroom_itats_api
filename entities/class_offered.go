package entities

type ClassOffered struct {
	SubjectClass            string  `gorm:"column:kelas" json:"subject_class"`
	SubjectCredits          int16   `gorm:"column:mksks" json:"subject_credit"`
	SubjectID               string  `gorm:"column:mkid" json:"subject_id"`
	SubjectName             string  `gorm:"column:mknama" json:"subject_name"`
	MajorID                 string  `gorm:"column:jurid" json:"major_id"`
	AcademicPeriodID        string  `gorm:"column:pakid" json:"academic_period_id"`
	LecturerID              string  `gorm:"column:dosid" json:"lecturer_id"`
	ActivityMasterID        string  `gorm:"column:id_master_kegiatan" json:"activity_master_id"`
	PresencePercentageScore float64 `gorm:"column:nilai_persentase_kehadiran" json:"presence_percentage_score"`
}
