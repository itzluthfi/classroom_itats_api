package student_repositories

import (
	"classroom_itats_api/entities"
	"context"
	"time"

	"gorm.io/gorm"
)

type studentPresenceRepository struct {
	db *gorm.DB
}

type StudentPresenceRepository interface {
	GetStudentPresences(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Presence, error)
	GetPresenceQuestion(ctx context.Context, pakID string) ([]entities.PresenceQuestion, error)
	GetSubjectResponsi(ctx context.Context, pakID string, mkID string, class string) (int, error)
	SetStudentPresenceAnswers(ctx context.Context, PresenceAnswers []entities.PresenceAnswer) error
	SetStudentPresence(ctx context.Context, StudentPresence entities.Presence) error
	PresenceCreated(ctx context.Context) ([]map[string]interface{}, error)
	PresenceReminder(ctx context.Context) ([]map[string]interface{}, error)
}

func NewStudentPresenceRepository(db *gorm.DB) *studentPresenceRepository {
	return &studentPresenceRepository{
		db: db,
	}
}

func (s *studentPresenceRepository) GetStudentPresences(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Presence, error) {
	studentPresence := []entities.Presence{}

	err := s.db.WithContext(ctx).Table("absen").Where("pakid = ?", pakID).Where("mkid = ?", mkID).Where("kelas = ?", class).Where("mhsid = ?", mhsID).Order("kultype ASC, weekid ASC").Find(&studentPresence).Error

	return studentPresence, err
}

func (s *studentPresenceRepository) GetPresenceQuestion(ctx context.Context, pakID string) ([]entities.PresenceQuestion, error) {
	presenceQuestions := []entities.PresenceQuestion{}

	err := s.db.WithContext(ctx).Table("master_pertanyaan").Joins("JOIN master_pertanyaan_pak on master_pertanyaan.master_pertanyaan_id = master_pertanyaan_pak.master_pertanyaan_id").Where("master_pertanyaan_pak.pakid = ?", pakID).Find(&presenceQuestions).Error

	return presenceQuestions, err
}

func (s *studentPresenceRepository) GetSubjectResponsi(ctx context.Context, pakID string, mkID string, class string) (int, error) {
	var rows int

	err := s.db.WithContext(ctx).Raw(`select count(*) from vw_jadwal where pakid = ? and mkid = ? and kelas = ? and kultipeid != ?`, pakID, mkID, class, "M").Find(&rows).Error

	return rows, err
}

func (s *studentPresenceRepository) SetStudentPresenceAnswers(ctx context.Context, PresenceAnswers []entities.PresenceAnswer) error {
	return s.db.WithContext(ctx).Table("kul_jawaban").Create(PresenceAnswers).Error
}

func (s *studentPresenceRepository) SetStudentPresence(ctx context.Context, StudentPresence entities.Presence) error {
	return s.db.WithContext(ctx).Table("absen").Create(StudentPresence).Error
}

func (s *studentPresenceRepository) PresenceCreated(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	kuls := []entities.Lecture{}

	err := s.db.WithContext(ctx).Table("kul").Where("waktu_entri between ? AND ?", time.Now().Add(time.Duration(-30)*time.Minute), time.Now()).Where("batas_presensi > ?", time.Now()).Find(&kuls).Error
	if err != nil {
		return nil, err
	}

	for _, kul := range kuls {
		var subject string
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string
		e := s.db.WithContext(ctx).Table("mk").Select("mknama").Where("mkid = ?", kul.SubjectID).Find(&subject).Error
		if e != nil {
			err = e
		}

		e = s.db.WithContext(ctx).Table("jur").Select("jurid").Where("jurid = ?", kul.MajorID).Or("jur_parent_id = ?", kul.MajorID).Find(&jurID).Error

		if e != nil {
			err = e
		}

		if len(jurID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("krs").Select("mhsid").
				Where("pakid = ?", kul.AcademicPeriodID).
				Where("mkid = ?", kul.SubjectID).
				Where("kelaskrs = ?", kul.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT mhsid FROM absen where absen.pakid = krs.pakid and absen.mkid = krs.mkid and absen.kelas = krs.kelaskrs and absen.mhsid = krs.mhsid and absen.weekid = ? and absen.kultype = ? and absen.jurid = ?)", kul.WeekID, kul.LectureType, kul.MajorID).Find(&mhsID).Error
		}

		if e != nil {
			err = e
		}

		if len(mhsID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("users").Select("mobile_token").Where("name in ?", mhsID).Where("mobile_token != ?", "null").Find(&user).Error
		}

		if e == nil {
			usr["kul"] = kul
			usr["user"] = user
			usr["subject"] = subject
			users = append(users, usr)
		}

		if err != nil {
			err = e
		}

	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *studentPresenceRepository) PresenceReminder(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	kuls := []entities.Lecture{}

	err := s.db.WithContext(ctx).Table("kul").Where("batas_presensi between ? and ?", time.Now(), time.Now().Add(time.Duration(3)*time.Hour)).Find(&kuls).Error
	if err != nil {
		return nil, err
	}

	for _, kul := range kuls {
		var subject string
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string
		e := s.db.WithContext(ctx).Table("mk").Select("mknama").Where("mkid = ?", kul.SubjectID).Find(&subject).Error
		if e != nil {
			err = e
		}

		e = s.db.WithContext(ctx).Table("jur").Select("jurid").Where("jurid = ?", kul.MajorID).Or("jur_parent_id = ?", kul.MajorID).Find(&jurID).Error

		if e != nil {
			err = e
		}

		if len(jurID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", kul.AcademicPeriodID).
				Where("mkid = ?", kul.SubjectID).
				Where("kelaskrs = ?", kul.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT mhsid FROM absen where absen.pakid = krs.pakid and absen.mkid = krs.mkid and absen.kelas = krs.kelaskrs and absen.mhsid = krs.mhsid and absen.weekid = ? and absen.kultype = ? and absen.jurid = ?)", kul.WeekID, kul.LectureType, kul.MajorID).Find(&mhsID).Error
		}

		if e != nil {
			err = e
		}

		if len(mhsID) > 0 || e == nil {
			err = s.db.WithContext(ctx).Table("users").Select("mobile_token").Where("name in ?", mhsID).Where("mobile_token != ?", "null").Find(&user).Error
		}

		if e == nil {
			usr["kul"] = kul
			usr["user"] = user
			usr["subject"] = subject
			users = append(users, usr)
		}

		if err != nil {
			err = e
		}

	}

	if err != nil {
		return nil, err
	}

	return users, nil
}
