package input

type LecturerSubjectFilter struct {
	AcademicPeriodID string `json:"period" binding:"required"`
	MajorID          string `json:"major" binding:"-"`
}

type LecturerMajor struct {
	AcademicPeriodID string `json:"period" binding:"-"`
}
