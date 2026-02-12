package entities

type ScoreType struct {
	ScoreTypeID   string `gorm:"column:jnilid" json:"score_type_id"`
	ScoreTypeDesc string `gorm:"column:jnildesc" json:"score_type_desc"`
	ScoreWeight   int16  `gorm:"column:bobot" json:"score_weight"`
	IsActive      bool   `gorm:"column:jnilaktif" json:"is_active"`
	MinimumScore  int    `gorm:"column:batas_nilai_awal" json:"minimum_score"`
	MaximumScore  int    `gorm:"column:batas_nilai_akhir" json:"maximum_score"`
}
