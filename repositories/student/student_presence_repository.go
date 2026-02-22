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
	GetActivePresence(ctx context.Context, mkID string, pakID string, class string, mhsID string, dosID string) ([]map[string]interface{}, error)
	GetHomeActivePresence(ctx context.Context, mhsID string) ([]map[string]interface{}, error)
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

func (s *studentPresenceRepository) GetActivePresence(ctx context.Context, mkID string, pakID string, class string, mhsID string, dosID string) ([]map[string]interface{}, error) {
	presences := []entities.Lecture{}
	result := []map[string]interface{}{}

	todayDate := time.Now().Format("2006-01-02")
	now := time.Now()

	err := s.db.WithContext(ctx).Table("kul").
		Select("kul.*, jam.jammulai, jam.jamhingga, kultipe.kultipenama").
		Joins("left join jam on kul.jamid = jam.jamid").
		Joins("left join kultipe on kul.kultype = kultipe.kultipeid").
		Where("kul.mkid = ?", mkID).
		Where("kul.pakid = ?", pakID).
		Where("kul.kelas = ?", class).
		Where("kul.dosid = ?", dosID).
		Where("kul.kultgl = ?", todayDate).
		Where("kul.is_kelas_hadir = ?", true).
		Where("kul.batas_presensi IS NOT NULL").
		Where("kul.batas_presensi >= ?", now).
		Order("jam.jammulai asc").
		Find(&presences).Error

	if err != nil {
		return nil, err
	}

	for _, kul := range presences {
		var count int64
		s.db.WithContext(ctx).Table("absen").
			Where("pakid = ?", kul.AcademicPeriodID).
			Where("mkid = ?", kul.SubjectID).
			Where("kelas = ?", kul.SubjectClass).
			Where("kultgl = ?", kul.LectureSchedule).
			Where("kultype = ?", kul.LectureType).
			Where("jamid = ?", kul.HourID).
			Where("weekid = ?", kul.WeekID).
			Where("mhsid = ?", mhsID).
			Count(&count)

		sudahPresensi := count > 0

		item := map[string]interface{}{
			"kul":            kul,
			"sudah_presensi": sudahPresensi,
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *studentPresenceRepository) GetHomeActivePresence(ctx context.Context, mhsID string) ([]map[string]interface{}, error) {
	presences := []entities.Lecture{}
	result := []map[string]interface{}{}

	todayDate := time.Now().Format("2006-01-02")
	now := time.Now()

	// Get active pakid first
	var pakID string
	err := s.db.WithContext(ctx).Table("(?) as pak", s.db.Table("pak").Order("pakid DESC").Limit(5)).
		Select("pakid").Where("isactive = ?", true).Find(&pakID).Error
	if err != nil {
		return nil, err
	}

	// Get all active presences today across all subjects the student is enrolled in
	err = s.db.WithContext(ctx).Table("kul").
		Select("kul.*, jam.jammulai, jam.jamhingga, kultipe.kultipenama, mk.mknama").
		Joins("left join jam on kul.jamid = jam.jamid").
		Joins("left join kultipe on kul.kultype = kultipe.kultipeid").
		Joins("left join mk on kul.mkid = mk.mkid").
		Joins("join krs on krs.mkid = kul.mkid and krs.kelaskrs = kul.kelas and krs.pakid = kul.pakid").
		Where("krs.mhsid = ?", mhsID).
		Where("kul.pakid = ?", pakID).
		Where("kul.kultgl = ?", todayDate).
		Where("kul.is_kelas_hadir = ?", true).
		Where("kul.batas_presensi IS NOT NULL").
		Where("kul.batas_presensi >= ?", now).
		Order("jam.jammulai asc").
		Find(&presences).Error

	if err != nil {
		return nil, err
	}

	for _, kul := range presences {
		var count int64
		s.db.WithContext(ctx).Table("absen").
			Where("pakid = ?", kul.AcademicPeriodID).
			Where("mkid = ?", kul.SubjectID).
			Where("kelas = ?", kul.SubjectClass).
			Where("kultgl = ?", kul.LectureSchedule).
			Where("kultype = ?", kul.LectureType).
			Where("jamid = ?", kul.HourID).
			Where("weekid = ?", kul.WeekID).
			Where("mhsid = ?", mhsID).
			Count(&count)

		sudahPresensi := count > 0

		item := map[string]interface{}{
			"kul":            kul,
			"sudah_presensi": sudahPresensi,
		}
		result = append(result, item)
	}

	return result, nil
}
