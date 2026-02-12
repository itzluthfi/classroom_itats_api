package entities

type SubjectJSON struct {
	SubjectClass     string            `json:"subject_class"`
	SubjectCredits   int16             `json:"subject_credit"`
	SubjectID        string            `json:"subject_id"`
	MajorID          string            `json:"major_id"`
	MajorName        string            `json:"major_name"`
	AcademicPeriodID string            `json:"academic_period_id"`
	LecturerName     string            `json:"lecturer_name"`
	LecturerID       string            `json:"lecturer_id"`
	SubjectName      string            `json:"subject_name"`
	TotalStudent     int               `json:"total_student"`
	ActivityMasterID string            `json:"activity_master_id"`
	SubjectSchedules []SubjectSchedule `json:"subject_schedules"`
}

type StudentSubject struct {
	SubjectClass     string `gorm:"column:kelaskrs"`
	SubjectCredits   int16  `gorm:"column:sks"`
	SubjectID        string `gorm:"column:mkid"`
	MajorID          string `gorm:"column:jurid"`
	AcademicPeriodID string `gorm:"column:pakid"`
	LecturerName     string `gorm:"column:dosnama"`
	LecturerID       string `gorm:"column:dosid"`
	SubjectName      string `gorm:"column:mknama"`
	ActivityMasterID string `gorm:"column:id_master_kegiatan"`
}

type LecturerSubject struct {
	SubjectClass     string `gorm:"column:kelas" json:"subject_class"`
	SubjectCredits   int16  `gorm:"column:mksks" json:"subject_credit"`
	SubjectID        string `gorm:"column:mkid" json:"subject_id"`
	MajorID          string `gorm:"column:jurid" json:"major_id"`
	MajorName        string `gorm:"column:jurasli_jurnama" json:"major_name"`
	AcademicPeriodID string `gorm:"column:pakid" json:"academic_period_id"`
	LecturerID       string `gorm:"column:dosid" json:"lecturer_id"`
	LecturerName     string `gorm:"column:dosnama" json:"lecturer_name"`
	SubjectName      string `gorm:"column:mknama" json:"subject_name"`
	TotalStudent     int    `gorm:"column:jumlah_mhs_perkelas" json:"total_student"`
	ActivityMasterID string `gorm:"column:id_master_kegiatan" json:"activity_master_id"`
}

type LecturerSubjectReport struct {
	SubjectClass     string `gorm:"column:kelas" json:"subject_class"`
	SubjectCredits   int16  `gorm:"column:mksks" json:"subject_credit"`
	SubjectID        string `gorm:"column:mkid" json:"subject_id"`
	MajorID          string `gorm:"column:jurid" json:"major_id"`
	MajorName        string `gorm:"column:jurasli_jurnama" json:"major_name"`
	AcademicPeriodID string `gorm:"column:pakid" json:"academic_period_id"`
	LecturerID       string `gorm:"column:dosid" json:"lecturer_id"`
	LecturerName     string `gorm:"column:dosnama" json:"lecturer_name"`
	SubjectName      string `gorm:"column:mknama" json:"subject_name"`
	TotalStudent     int    `gorm:"column:jumlah_mhs_perkelas" json:"total_student"`
	SubjectType      string `gorm:"column:kultipeid" json:"subject_type"`
	SubjectRoom      string `gorm:"column:ruangid" json:"subject_room"`
	Day              string `gorm:"column:haridesc" json:"day"`
	TimeStart        string `gorm:"column:jammulai" json:"time_start"`
	TimeEnd          string `gorm:"column:jamhingga" json:"time_end"`
	HourID           string `gorm:"column:jamid" json:"hour_id"`
	ActivityMasterID string `gorm:"column:id_master_kegiatan" json:"activity_master_id"`
}

type SubjectSchedule struct {
	SubjectID    string `gorm:"column:mkid" json:"subject_id"`
	SubjectClass string `gorm:"column:kelas" json:"subject_class"`
	SubjectType  string `gorm:"column:kultipeid" json:"subject_type"`
	SubjectRoom  string `gorm:"column:ruangid" json:"subject_room"`
	Day          string `gorm:"column:haridesc" json:"day"`
	TimeStart    string `gorm:"column:jammulai" json:"time_start"`
	TimeEnd      string `gorm:"column:jamhingga" json:"time_end"`
}
