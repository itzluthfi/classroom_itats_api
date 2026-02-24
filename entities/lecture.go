package entities

import "time"

type Lecture struct {
	AcademicPeriodID    string    `gorm:"column:pakid" json:"academic_period_id"`
	SubjectID           string    `gorm:"column:mkid" json:"subject_id"`
	MajorID             string    `gorm:"column:jurid" json:"major_id"`
	LecturerID          string    `gorm:"column:dosid" json:"lecturer_id"`
	SubjectClass        string    `gorm:"column:kelas" json:"subject_class"`
	LectureSchedule     time.Time `gorm:"column:kultgl" json:"lecture_schedule"`
	LectureType         string    `gorm:"column:kultype" json:"lecture_type"`
	SubjectCredits      int16     `gorm:"column:sks" json:"subject_credit"`
	HourID              string    `gorm:"column:jamid" json:"hour_id"`
	Material            string    `gorm:"column:materi" json:"material_realization"`
	LectureLink         string    `gorm:"column:link_kuliah" json:"lecture_link"`
	EntryTime           time.Time `gorm:"column:waktu_entri" json:"entry_time"`
	ApprovalStatus      int16     `gorm:"column:status_approval" json:"approval_status"`
	WeekID              int16     `gorm:"column:weekid" json:"week_id"`
	TimeRealization     int16     `gorm:"column:realisasiwaktu" json:"time_realization"`
	TimeSuitability     bool      `gorm:"column:kesesuaianwaktu" json:"time_suitability"`
	MaterialSuitability bool      `gorm:"column:kesesuaianmateri" json:"material_suitability"`
	MaterialLink        string    `gorm:"column:link_materi" json:"material_link"`
	LectureID           string    `gorm:"column:kulid" json:"lecture_id"`
	PresenceLimit       time.Time `gorm:"column:batas_presensi" json:"presence_limit"`
	PresenceStudent     int       `gorm:"column:mahasiswa_absen" json:"presence_student"`
	LinkMeet            string    `gorm:"column:link_meet" json:"link_meet"`
	LinkRecord          string    `gorm:"column:link_record" json:"link_record"`
	CollegeType         int       `gorm:"column:jenis_kuliah" json:"college_type"`
	CollegeTypeName     string    `gorm:"column:nama_jenis_perkuliahan" json:"college_type_name"`
	TimeStart           string    `gorm:"column:jammulai" json:"time_start"`
	TimeEnd             string    `gorm:"column:jamhingga" json:"time_end"`
	LectureTypeName     string    `gorm:"column:kultipenama" json:"lecture_type_name"`
	SubjectName         string    `gorm:"column:mknama" json:"subject_name"`
}

type LectureWeek struct {
	AcademicPeriodID string `gorm:"column:pakid" json:"academic_period_id"`
	SubjectID        string `gorm:"column:mkid" json:"subject_id"`
	MajorID          string `gorm:"column:jurid" json:"major_id"`
	LecturerID       string `gorm:"column:dosid" json:"lecturer_id"`
	SubjectClass     string `gorm:"column:kelas" json:"subject_class"`
	LectureType      string `gorm:"column:kultype" json:"lecture_type"`
	HourID           string `gorm:"column:jamid" json:"hour_id"`
	WeekID           int16  `gorm:"column:weekid" json:"week_id"`
	LectureID        string `gorm:"column:kulid" json:"lecture_id"`
	LinkMeet         string `gorm:"column:link_meet" json:"link_meet"`
	LinkRecord       string `gorm:"column:link_record" json:"link_record"`
	CollegeType      int    `gorm:"column:jenis_kuliah" json:"college_type"`
	CollegeTypeName  string `gorm:"column:nama_jenis_perkuliahan" json:"college_type_name"`
}

type StudentLectureJSON struct {
	MaterialLectures []Lecture   `json:"material_lectures"`
	ResponsiLectures [][]Lecture `json:"responsi_lectures"`
}

type StoreLecture struct {
	LectureID           string    `gorm:"column:kulid" json:"lecture_id"`
	LecturerID          string    `gorm:"column:dosid" json:"lecturer_id"`
	MajorID             string    `gorm:"column:jurid" json:"major_id"`
	SubjectID           string    `gorm:"column:mkid" json:"subject_id"`
	SubjectClass        string    `gorm:"column:kelas" json:"subject_class"`
	HourID              string    `gorm:"column:jamid" json:"hour_id"`
	LectureType         string    `gorm:"column:kultype" json:"lecture_type"`
	TimeRealization     int16     `gorm:"column:realisasiwaktu" json:"time_realization"`
	WeekID              int16     `gorm:"column:weekid" json:"week_id"`
	SubjectCredits      int16     `gorm:"column:sks" json:"subject_credit"`
	AcademicPeriodID    string    `gorm:"column:pakid" json:"academic_period_id"`
	LectureSchedule     time.Time `gorm:"column:kultgl" json:"lecture_schedule"`
	ApprovalStatus      int16     `gorm:"column:status_approval" json:"approval_status"`
	EntryTime           time.Time `gorm:"column:waktu_entri" json:"entry_time"`
	MaterialRealization string    `gorm:"column:materi" json:"material_realization"`
	PresenceLimit       time.Time `gorm:"column:batas_presensi" json:"presence_limit"`
	CollegeType         int       `gorm:"column:jenis_kuliah" json:"college_type"`
}

type LectureMaterial struct {
	LectureMaterialID string `gorm:"column:kul_materi_id"`
	LectID            string `gorm:"column:kul_id"`
	MaterialID        string `gorm:"column:materi_id" json:"material_id"`
}

type LectureRequest struct {
	LectureStore     StoreLecture      `json:"lecture"`
	LectureMaterials []LectureMaterial `json:"material"`
}

type LectureEditRequest struct {
	LectureStore     Lecture           `json:"lecture"`
	LectureMaterials []LectureMaterial `json:"material"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (LectureMaterial) TableName() string {
	return "kul_materi"
}

func (StoreLecture) TableName() string {
	return "kul"
}
