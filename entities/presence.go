package entities

import "time"

type StudentPresence struct {
	PresenceStudent Presence         `json:"presence_student"`
	PresenceAnswers []PresenceAnswer `json:"presence_answers"`
}

type Presence struct {
	AcademicPeriodID string    `gorm:"column:pakid" json:"academic_period_id"`
	SubjectID        string    `gorm:"column:mkid" json:"subject_id"`
	MajorID          string    `gorm:"column:jurid" json:"major_id"`
	SubjectClass     string    `gorm:"column:kelas" json:"subject_class"`
	CollegeSchedule  time.Time `gorm:"column:kultgl" json:"college_schedule"`
	StudentID        string    `gorm:"column:mhsid" json:"student_id"`
	IsPresent        bool      `gorm:"column:ishadir" json:"is_present"`
	CollegeType      string    `gorm:"column:kultype" json:"college_type"`
	HourID           string    `gorm:"column:jamid" json:"hour_id"`
	WeekID           int16     `gorm:"column:weekid" json:"week_id"`
	IsOffline        bool      `gorm:"column:isoffline" json:"is_offline"`
	Score            int16     `gorm:"column:nilai" json:"score"`
}

type PresenceQuestion struct {
	MasterQuestionID int       `gorm:"column:master_pertanyaan_id" json:"master_question_id"`
	Question         string    `gorm:"column:pertanyaan" json:"question"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type PresenceQuestionPeriod struct {
	MasterQuestionPeriodID int    `gorm:"column:mpp_id" json:"master_question_period_id"`
	PeriodID               string `gorm:"column:pakid" json:"period_id"`
	MasterQuestionID       int    `gorm:"column:master_pertanyaan_id" json:"master_question_id"`
}

type PresenceAnswer struct {
	CollegeID        string    `gorm:"column:kulid;type:uuid" json:"lecture_id"`
	StudentID        string    `gorm:"column:mhsid"`
	MasterQuestionID int       `gorm:"column:mp_id" json:"presence_question_id"`
	Answer           int       `gorm:"column:jawaban" json:"answer"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type StudentPresenceJSON struct {
	MaterialPresences []Presence   `json:"material_presences"`
	ResponsiPresences [][]Presence `json:"responsi_presences"`
}
