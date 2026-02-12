package entities

type SubjectMember struct {
	UserID      string `gorm:"column:uid" json:"user_id"`
	Name        string `gorm:"column:nama" json:"name"`
	PhoneNumber string `gorm:"column:mobile" json:"phone_number"`
	Presence    int    `gorm:"column:kehadiran" json:"presence"`
}

type StudentProfileJSON struct {
	UserID                  string                   `json:"user_id"`
	Name                    string                   `json:"name"`
	Photo                   string                   `json:"photo"`
	PhoneNumber             string                   `json:"phone_number"`
	Email                   string                   `json:"email"`
	Presence                int                      `json:"presence"`
	TotalPresence           int                      `json:"total_presence"`
	AssignmentSubmited      int                      `json:"assignment_submited"`
	TotalAssignment         int                      `json:"total_assignment"`
	StudentSubjectPresences []StudentSubjectPresence `json:"student_subject_presences"`
}

type StudentProfile struct {
	UserID             string `gorm:"column:mhsid"`
	Name               string `gorm:"column:mhsnama"`
	Photo              string `gorm:"column:foto"`
	PhoneNumber        string `gorm:"column:mobile"`
	Email              string `gorm:"column:email"`
	Presence           int    `gorm:"column:absen_mhs"`
	TotalPresence      int    `gorm:"column:total_absen"`
	AssignmentSubmited int    `gorm:"column:tugas_terkumpul"`
	TotalAssignment    int    `gorm:"column:total_tugas"`
}

type StudentSubjectPresence struct {
	SubjectID        string `gorm:"column:mkid" json:"subject_id"`
	SubjectClass     string `gorm:"column:kelas" json:"subject_class"`
	SubjectName      string `gorm:"column:mknama" json:"subject_name"`
	ActivityMasterID string `gorm:"column:id_master_kegiatan" json:"activity_master_id"`
	Presence         int    `gorm:"column:absen_mhs" json:"presence"`
	TotalPresence    int    `gorm:"column:total_absen" json:"total_presence"`
}

type StudentCheck struct {
	UserID string `gorm:"column:mhsid" json:"npm"`
	Name   string `gorm:"column:mhsnama" json:"nama"`
}

type StudentGetAll struct {
	UserID      string `gorm:"column:mhsid" json:"npm"`
	Name        string `gorm:"column:mhsnama" json:"nama"`
	PhoneNumber string `gorm:"column:mobile" json:"no_hp"`
	Email       string `gorm:"column:email" json:"email"`
	BirthDate   string `gorm:"column:tgllahir" json:"tgl_lahir"`
	Gender      string `gorm:"column:sexid" json:"gender"`
}

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// type Krs struct {
// 	MhsID     string `gorm:"type:varchar(20);not null"`
// 	JurID     string `gorm:"type:varchar(10);not null"`
// 	PakID     string `gorm:"type:varchar(10);not null"`
// 	MkID      string `gorm:"type:varchar(10);not null"`
// 	KelasKrs  string `gorm:"type:varchar(5);not null"`
// 	SKS       int16  `gorm:"type:smallint"`
// 	Jur       Jur    `gorm:"foreignkey:JurID;association_foreignkey:JurID"` // 1
// 	Mk        Mk     `gorm:"foreignkey:MkID;association_foreignkey:MkID"`   // 1
// 	Pak       Pak    `gorm:"foreignkey:PakID;association_foreignkey:PakID"` // 1
// 	Mahasiswa Mhs    `gorm:"foreignkey:MhsID;association_foreignkey:MhsID"` // 1
// }

// type Mhs struct {
// 	MhsID    string    `gorm:"type:varchar(20);primary_key"`
// 	MhsNama  string    `gorm:"type:varchar(50)"`
// 	TglLahir time.Time `gorm:"type:date"`
// 	Tlp      string    `gorm:"type:varchar(50)"`
// 	Mobile   string    `gorm:"type:varchar(50)"`
// 	Email    string    `gorm:"type:varchar(50)"`
// 	Aktif    string    `gorm:"type:varchar(1)"`
// 	UsrID    string    `gorm:"type:varchar(15)"`
// 	Foto     string    `gorm:"type:text"`
// 	AktifMhs Aktif     `gorm:"foreignKey:Aktif"`
// 	User     User      `gorm:"foreignKey:User"`
// }

// type Jur struct {
// 	JurID       string `gorm:"type:varchar(10);not null;primary_key"`
// 	JurNama     string `gorm:"type:varchar(50)"`
// 	JurNume     string `gorm:"type:varchar(50)"`
// 	ProdiNama   string `gorm:"type:varchar(40)"`
// 	ProdiName   string `gorm:"type:varchar(40)"`
// 	JurParentID string `gorm:"type:varchar(10)"`
// }

// type Mk struct {
// 	MkId         string `gorm:"type:varchar(10);primaryKey"`
// 	ItemId       string `gorm:"type:varchar(10)"`
// 	JurId        string `gorm:"type:varchar(10)"`
// 	MkNama       string `gorm:"type:varchar(50)"`
// 	MkDesc       string `gorm:"type:varchar(250)"`
// 	IsActive     bool
// 	IdMataKuliah string `gorm:"type:uuid"`
// 	MkSks        int16
// 	Jur          Jur `gorm:"foreignKey:JurId"` // 1
// }

// type Dos struct {
// 	DosID         string                 `gorm:"type:varchar(20);not null;primary_key"`
// 	NIDN          string                 `gorm:"type:varchar(20)"`
// 	DosNama       string                 `gorm:"type:varchar(50);not null"`
// 	Gelar         string                 `gorm:"type:varchar(20)"`
// 	Tlp           string                 `gorm:"type:varchar(50)"`
// 	Mobile        string                 `gorm:"type:varchar(50)"`
// 	Email         string                 `gorm:"type:varchar(50)"`
// 	Photo         string                 `gorm:"type:varchar(255)"`
// 	GelarBaru     map[string]interface{} `gorm:"type:jsonb;serializer:json"`
// 	Aktif         bool                   `gorm:"not null;default:true"`
// 	DosKategoriID int
// 	Foto          string      `gorm:"type:varchar(100)"`
// 	DosKategori   DosKategori `gorm:"foreignkey:doskategoriid"`
// 	Usr           User        `gorm:"foreignkey:usrid"`
// }

// // type Jad struct {
// // 	JamId            string `gorm:"type:varchar(5);primaryKey"`
// // 	PakId            string `gorm:"type:varchar(10);primaryKey"`
// // 	MkId             string `gorm:"type:varchar(10);primaryKey"`
// // 	JurId            string `gorm:"type:varchar(10);primaryKey"`
// // 	Kelas            string `gorm:"type:varchar(5);default:'A';primaryKey"`
// // 	JadNo            int16  `gorm:"primaryKey"`
// // 	Hari             string `gorm:"type:varchar(1)"`
// // 	RuangId          string `gorm:"type:varchar(6)"`
// // 	JadDesc          string `gorm:"type:varchar(50)"`
// // 	Sks              int16
// // 	Responsi         bool                   `gorm:"default:false"`
// // 	Praktikum        bool                   `gorm:"default:false"`
// // 	DosId            string                 `gorm:"type:varchar(20)"`
// // 	KulTipeId        string                 `gorm:"type:varchar(15);primaryKey"`
// // 	IdMasterKegiatan string                 `gorm:"type:uuid"`
// // 	JadWeek          map[string]interface{} `gorm:"type:jsonb;serializer:json"`
// // 	JadTeam          map[string]interface{} `gorm:"type:jsonb;serializer:json"`
// // 	HariJad          Hari                   `gorm:"foreignKey:Hari"`
// // 	Jam              Jam                    `gorm:"foreignKey:Jam"`
// // 	Ruang            Ruang                  `gorm:"foreignKey:RuangId"`
// // 	Dos              Dos                    `gorm:"foreignKey:DosId"`
// // 	KulTipe          KulTipe                `gorm:"foreignKey:KulTipeId"`
// // }

// type Klstw struct {
// 	JurID     string `gorm:"type:varchar(10);not null;primaryKey"`
// 	PakID     string `gorm:"type:varchar(10);not null;primaryKey"`
// 	MkID      string `gorm:"type:varchar(10);not null;primaryKey"`
// 	Kelas     string `gorm:"type:varchar(5);not null;default:'A';primaryKey"`
// 	DosID     string `gorm:"type:varchar(20)"`
// 	KlstwKet  string `gorm:"type:varchar(50)"`
// 	Kapasitas int    `gorm:"type:integer;default:60"`
// 	SemTW     int16
// 	Jur       Jur `gorm:"foreignKey:JurID"`
// 	Pak       Pak `gorm:"foreignKey:PakID"`
// 	Mk        Mk  `gorm:"foreignKey:MkID"`
// 	Dos       Dos `gorm:"foreignKey:DosID"`
// }

// type Ruang struct {
// 	RuangId string `gorm:"type:varchar(6);primaryKey"`
// }

// type KulTipe struct {
// 	KulTipeId string `gorm:"type:varchar(15);primaryKey"`
// }

// type Hari struct {
// 	Hari string `gorm:"type:varchar(1);primaryKey"`
// }

// type Jam struct {
// 	JamId string `gorm:"type:varchar(5);primaryKey"`
// }

// type DosKategori struct {
// 	gorm.Model
// 	DosKategoriID int `gorm:"type:int"`
// }

// type Prop struct {
// 	PropID string `gorm:"type:varchar(3);primaryKey"`
// 	// ... other fields ...
// }

// type Pek struct {
// 	PekID string `gorm:"type:varchar(3);primaryKey"`
// 	// ... other fields ...
// }

// type Pak struct {
// 	PakID string `gorm:"type:varchar(10);primaryKey"`
// 	// ... other fields ...
// }

// type Phs struct {
// 	PhsID string `gorm:"type:varchar(2);primaryKey"`
// 	// ... other fields ...
// }

// type Aktif struct {
// 	AktifID string `gorm:"type:varchar(1);primaryKey"`
// 	// ... other fields ...
// }

// type Item struct {
// 	ItemId string `gorm:"type:varchar(10);primaryKey"`
// }

// type ProgramStudi struct {
// }
