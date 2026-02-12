package entities

type PercentageScore struct {
	TotalPercentage        float64                 `json:"total_percentage"`
	PercentageScoreDetails []PercentageScoreDetail `json:"percentage_score_detail"`
}

type PercentageScoreDetail struct {
	ID              string  `json:"id"`
	AssignmentTitle string  `json:"assignment_title"`
	AssignmentType  string  `json:"assignment_type"`
	WeekID          string  `json:"week_id"`
	Percentage      float64 `json:"percentage"`
}
