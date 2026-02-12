package repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type lectureRepository struct {
	db *gorm.DB
}

type LectureRepository interface {
	GetLecture(ctx context.Context, pakID string, subjectID string, class string) ([]entities.Lecture, error)
}

func NewLectureRepository(db *gorm.DB) *lectureRepository {
	return &lectureRepository{
		db: db,
	}
}

func (l *lectureRepository) GetLecture(ctx context.Context, pakID string, subjectID string, class string) ([]entities.Lecture, error) {
	lectures := []entities.Lecture{}

	// err := l.db.WithContext(ctx).Raw(`select kul.*, count(absen.*) as mahasiswa_absen from kul join absen on absen.weekid = kul.weekid and absen.jamid = kul.jamid
	// 		and absen.jurid = kul.jurid and absen.kelas = kul.kelas and absen.kultgl = kul.kultgl and absen.mkid = kul.mkid
	// 		and absen.pakid = kul.pakid
	// 		where kul.mkid = ? and kul.kelas = ?
	// 		group by kul.jurid, kul.pakid, kul.batas_presensi, kul.dosid, kul.jamid, kul.kelas, kul.kesesuaianmateri,
	// 		kul.kesesuaianwaktu, kul.kulid, kul.kultgl, kul.kultype, kul.link_kuliah, kul.link_materi, kul.materi,
	// 		kul.mkid, kul.realisasiwaktu, kul.sks, kul.status_approval, kul.waktu_entri, kul.weekid
	// 		order by kul.weekid`, subjectID, class).Find(&lectures).Error

	err := l.db.WithContext(ctx).Raw(`select distinct kul.*, count(absen.*) as mahasiswa_absen from kul left join absen on absen.weekid = kul.weekid and absen.jamid = kul.jamid        
		and absen.jurid = kul.jurid and absen.kelas = kul.kelas and absen.kultgl = kul.kultgl and absen.mkid = kul.mkid
		and absen.pakid = kul.pakid and absen.kultype = kul.kultype
		where kul.pakid = ? and kul.mkid = ? and kul.kelas = ?
		group by kul.jurid, kul.pakid, kul.batas_presensi, kul.dosid, kul.jamid, kul.kelas,
		kul.kulid, kul.kultgl, kul.kultype, kul.link_kuliah, kul.link_materi, kul.materi,
		kul.mkid, kul.realisasiwaktu, kul.sks, kul.status_approval, kul.waktu_entri, kul.weekid
		order by kul.kultype, kul.weekid`, pakID, subjectID, class).Find(&lectures).Error

	return lectures, err
}
