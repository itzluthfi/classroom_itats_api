package entities

type StudentScore struct {
	StudentID       string `gorm:"column:mhsid" json:"student_id"`
	StudentName     string `gorm:"column:mhsnama" json:"student_name"`
	NumericScore    string `gorm:"column:nilagk" json:"numeric_score"`
	AlphabeticScore string `gorm:"column:nilhrf" json:"alphabetic_score"`
}

type Reg struct {
	RegID   string
	RegDesc string
	RegVal  string
}
