package entities

import "time"

type Assignment struct {
	AssignmentID     int       `gorm:"column:id_tugas_kul" json:"assignment_id"`
	ActivityMasterID string    `gorm:"column:master_kegiatan_id" json:"activity_master_id"`
	WeekID           int       `gorm:"column:weekid" json:"week_id"`
	AssignmentTitle  string    `gorm:"column:judul_tugas" json:"assignment_title"`
	Description      string    `gorm:"column:deskripsi" json:"description"`
	DueDate          time.Time `gorm:"column:batas_pengumpulan" json:"due_date"`
	StartTime        time.Time `gorm:"column:waktu_mulai_tugas" json:"start_time"`
	EndTime          time.Time `gorm:"column:waktu_akhir_tugas" json:"end_time"`
	JNilID           string    `gorm:"column:jnilid" json:"j_nil_id"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	FileLink         string    `gorm:"column:link_file" json:"file_link"`
	FileName         string    `gorm:"column:nama_file" json:"file_name"`
	IsShow           bool      `gorm:"column:is_tampil" json:"is_show"`
	RealPrercentage  float64   `gorm:"column:real_persentase" json:"real_prercentage"`
	SubjectClass     string    `gorm:"column:kelas" json:"subject_class"`
	Subjectname      string    `gorm:"column:mknama" json:"subject_name"`
	JnilDesc         string    `gorm:"column:jnildesc" json:"j_nil_desc"`
	TotalSubmited    int       `gorm:"column:jml_pengumpulan" json:"total_submited"`
}

type AssignmentJoin struct {
	AssignmentID           int       `gorm:"column:id_tugas_kul" json:"assignment_id"`
	ActivityMasterID       string    `gorm:"column:master_kegiatan_id" json:"activity_master_id"`
	WeekID                 int       `gorm:"column:weekid" json:"week_id"`
	AssignmentTitle        string    `gorm:"column:judul_tugas" json:"assignment_title"`
	Description            string    `gorm:"column:deskripsi" json:"description"`
	DueDate                time.Time `gorm:"column:batas_pengumpulan" json:"due_date"`
	JNilID                 string    `gorm:"column:jnilid" json:"j_nil_id"`
	FileLink               string    `gorm:"column:link_file" json:"file_link"`
	FileName               string    `gorm:"column:nama_file" json:"file_name"`
	AssignmentSubmissionID int       `gorm:"column:id_tugas_submission" json:"assignment_submission_id"`
	AssignmentFile         string    `gorm:"column:file_tugas" json:"assignment_file"`
	AssignmentLink         string    `gorm:"column:link_tugas" json:"assignment_link"`
	Note                   string    `gorm:"column:note" json:"note"`
	StudentID              string    `gorm:"column:mhsid" json:"student_id"`
	IDAssignment           int       `gorm:"column:tugas_kul_id" json:"id_assignment"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type AssignmentSubmission struct {
	AssignmentSubmissionID int       `gorm:"column:id_tugas_submission" json:"assignment_submission_id"`
	AssignmentFile         string    `gorm:"column:file_tugas" json:"assignment_file"`
	AssignmentLink         string    `gorm:"column:link_tugas" json:"assignment_link"`
	Note                   string    `gorm:"column:note" json:"note"`
	Score                  float64   `gorm:"column:nilai" json:"score"`
	StudentID              string    `gorm:"column:mhsid" json:"student_id"`
	AssignmentID           int       `gorm:"column:tugas_kul_id" json:"assignment_id"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
	FinalScore             float64   `gorm:"column:nilai_akhir" json:"final_score"`
}

type StudentAssignmentScore struct {
	AssignmentSubmissionID int     `gorm:"column:id_tugas_submission" json:"assignment_submission_id"`
	Score                  float64 `gorm:"column:nilai" json:"score"`
	StudentID              string  `gorm:"column:mhsid" json:"student_id"`
	FinalScore             float64 `gorm:"column:nilai_akhir" json:"final_score"`
	AssignmentID           int     `gorm:"column:tugas_kul_id" json:"assignment_id"`
	AssignmentTitle        string  `gorm:"column:judul_tugas" json:"assignment_title"`
}
