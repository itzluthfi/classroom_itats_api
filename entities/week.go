package entities

type Week struct {
	WeekID     int    `gorm:"column:weekid;type:integer" json:"week_id"`
	WeekNumber int    `gorm:"column:weekno;type:integer" json:"week_number"`
	IsActive   bool   `gorm:"column:isactive;type:boolean" json:"isactive"`
	IsTest     bool   `gorm:"column:isujian;type:boolean" json:"istest"`
	Note       string `gorm:"column:catatan;type:character varying(50)" json:"note"`
}
