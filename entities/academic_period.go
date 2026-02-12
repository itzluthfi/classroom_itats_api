package entities

type AcademicPeriod struct {
	AcademicPeriodID          string `gorm:"type:varchar(10);column:pakid;primaryKey" json:"academic_period_id"`
	OddEven                   string `gorm:"type:varchar(1);column:sgg" json:"odd_even"`
	CurriculumID              string `gorm:"type:varchar(5);column:krkid" json:"curriculum_id"`
	AcademicPeriodDescription string `gorm:"type:varchar(50);column:pakdesc" json:"academic_period_description"`
	YearStart                 int    `gorm:"type:int;column:th1" json:"year_start"`
	YearEnd                   int    `gorm:"type:int;column:th2" json:"year_end"`
	AcademicPeriodIndex       int    `gorm:"type:int;column:pakidx" json:"academic_period_index"`
	IsActive                  bool   `gorm:"type:boolean;column:isactive" json:"is_active"`
}
