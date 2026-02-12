package entities

type StudyAchievement struct {
	StudyAchievementID          int    `gorm:"column:id_capaian_pembelajaran" json:"study_achievement_id"`
	SubjectID                   string `gorm:"column:mkid" json:"subject_id"`
	StudyAchievementDescription string `gorm:"column:deskripsi_cp" json:"study_achievement_description"`
	LecturerID                  string `gorm:"column:dosid" json:"lecturer_id"`
	MajorPlanID                 int    `gorm:"column:capaian_jurusan_id" json:"major_plan_id"`
	StudyPlanID                 int    `gorm:"column:id_rencana_pembelajaran" json:"study_plan_id"`
	MajorID                     string `gorm:"column:jurid" json:"major_id"`
	AcademicPeriodID            string `gorm:"column:pakid" json:"academic_period_id"`
	SubjectClass                string `gorm:"column:kelas" json:"subject_class"`
	StudyPlanDescription        string `gorm:"column:deskripsi_rp" json:"study_plan_description"`
	Note                        string `gorm:"column:catatan" json:"note"`
	Code                        string `gorm:"column:kode" json:"code"`
	StudyPlanMappingID          int    `gorm:"column:id_mapping_rps" json:"study_plan_mapping_id"`
	WeekID                      int    `gorm:"column:weekid" json:"week_id"`
	LectureID                   string `gorm:"column:kulid" json:"lecture_id"`
}
