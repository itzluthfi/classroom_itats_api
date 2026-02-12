package input

type ForumInput struct {
	SubjectID    string `json:"subject_id" binding:"required"`
	SubjectClass string `json:"subject_class" binding:"required"`
}
