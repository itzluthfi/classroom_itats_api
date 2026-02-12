package entities

type User struct {
	UID        int    `gorm:"primarykey;column:uid;type:integer"`
	Name       string `gorm:"column:name;type:varchar(60)" json:"name"`
	Pass       string `gorm:"column:pass;type:varchar(60)" json:"pass"`
	Mail       string `gorm:"column:mail;type:varchar(60)" json:"mail"`
	Mode       int8   `gorm:"column:mode;type:smallint"`
	Sort       int8   `gorm:"column:sort;type:smallint"`
	Treshold   int8   `gorm:"column:treshold;type:smallint"`
	Theme      string `gorm:"column:theme;type:varchar(255)"`
	Signature  string `gorm:"column:signature;type:varchar(255)"`
	Created    int    `gorm:"column:created"`
	Access     int    `gorm:"column:access"`
	Login      int    `gorm:"column:login" json:"login"`
	Status     int8   `gorm:"column:status;type:smallint"`
	Timezone   string `gorm:"column:timezone;type:varchar(8)"`
	Language   string `gorm:"column:language;type:varchar(12)"`
	Picture    string `gorm:"column:picture;type:varchar(255)"`
	Init       string `gorm:"column:init;type:varchar(64)"`
	Data       string `gorm:"column:data"`
	Passwdasli string `gorm:"column:passwdasli;type:varchar(30)"`
}
