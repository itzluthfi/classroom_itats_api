package student_repositories

import (
	"classroom_itats_api/entities"
	"context"
	"strconv"

	"gorm.io/gorm"
)

type studentSubjectRepository struct {
	db *gorm.DB
}

type StudentSubjectRepository interface {
	GetActiveSubjectByStudentID(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error)
	GetSubjectClassStudent(ctx context.Context, mhsID string) ([]entities.SubjectSchedule, error)
	GetSubjectByStudentIDWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error)
	GetSubjectClassStudentWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectSchedule, error)
	StudyPeriodes(ctx context.Context, mhsID string) ([]entities.AcademicPeriod, error)
	GetActiveRoomLoans(ctx context.Context) ([]entities.RoomKeyLoan, error)
	GetTodayHariCode(ctx context.Context, dayOfWeek int) (string, error)
}

func NewStudentSubjectRepository(db *gorm.DB) *studentSubjectRepository {
	return &studentSubjectRepository{
		db: db,
	}
}

func (s *studentSubjectRepository) GetActiveSubjectByStudentID(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error) {
	studentSubjects := []entities.StudentSubject{}

	if pakID == "" {
		err := s.db.WithContext(ctx).Table("(?) as pak", s.db.Table("pak").Order("pakid DESC").Limit(5)).Select("pakid").Where("isactive = ?", true).Find(&pakID).Error

		if err != nil {
			return studentSubjects, err
		}
	}

	err := s.db.WithContext(ctx).Table("krs").
		Distinct("krs.mkid", "mk.mknama", "krs.kelaskrs", "krs.sks", "krs.mkid", "krs.jurid", "krs.pakid", "vw_kelas_tawar.dosid", "vw_kelas_tawar.dosnama", "vw_kelas_tawar.id_master_kegiatan").
		// Select("mk.mknama", "krs.kelaskrs", "krs.sks", "krs.mkid", "krs.jurid", "krs.pakid", "vw_kelas_tawar.dosid", "vw_kelas_tawar.dosnama", "vw_kelas_tawar.id_master_kegiatan").
		Joins("left join vw_kelas_tawar on vw_kelas_tawar.mkid = krs.mkid and vw_kelas_tawar.kelas = krs.kelaskrs and vw_kelas_tawar.pakid = krs.pakid").
		Joins("join mk on mk.mkid = krs.mkid").
		Joins("left join jad on jad.id_master_kegiatan = vw_kelas_tawar.id_master_kegiatan and mk.jurid = jad.jurid").
		// Joins("join pak on pak.pakid = krs.pakid").
		// Joins("join pakkul on pak.pakid = pakkul.pakid").
		// Joins("join jad on mk.mkid = jad.mkid and pak.pakid = jad.pakid and krs.kelaskrs = jad.kelas").
		Where("mhsid = ?", mhsID).
		// Where("pakkul.isactive = ?", true).
		Where("krs.pakid = ?", pakID).
		Order("mk.mknama ASC").
		Find(&studentSubjects).Error

	return studentSubjects, err
}

func (s *studentSubjectRepository) GetSubjectClassStudent(ctx context.Context, mhsID string) ([]entities.SubjectSchedule, error) {
	subjectSchedules := []entities.SubjectSchedule{}
	var pakID string

	err := s.db.WithContext(ctx).Table("(?) as pak", s.db.Table("pak").Order("pakid DESC").Limit(5)).Select("pakid").Where("isactive = ?", true).Find(&pakID).Error

	if err != nil {
		return []entities.SubjectSchedule{}, err
	}

	err = s.db.WithContext(ctx).Table("krs").
		Distinct("krs.mkid", "jad.kultipeid", "jad.ruangid", "hari.haridesc", "jad.hari as hari_code", "jam.jammulai", "jam.jamhingga", "jad.sks").
		Joins("join mk on mk.mkid = krs.mkid").
		Joins("join pak on pak.pakid = krs.pakid").
		Joins("join jad on mk.mkid = jad.mkid and pak.pakid = jad.pakid and krs.kelaskrs = jad.kelas").
		Joins("join jam on jam.jamid = jad.jamid").
		Joins("join hari on hari.hari = jad.hari").
		Where("mhsid = ?", mhsID).
		Where("pak.pakid = ?", pakID).
		Find(&subjectSchedules).Error

	if err != nil {
		return []entities.SubjectSchedule{}, err
	}

	return subjectSchedules, nil
}

func (s *studentSubjectRepository) GetSubjectByStudentIDWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error) {
	studentSubject := []entities.StudentSubject{}

	err := s.db.WithContext(ctx).Table("krs").
		Distinct("krs.mkid", "mk.mknama", "krs.kelaskrs", "krs.sks", "krs.mkid", "krs.jurid", "krs.pakid", "vw_kelas_tawar.dosid", "vw_kelas_tawar.dosnama", "vw_kelas_tawar.id_master_kegiatan").
		Joins("join vw_kelas_tawar on vw_kelas_tawar.mkid = krs.mkid and vw_kelas_tawar.kelas = krs.kelaskrs and vw_kelas_tawar.pakid = krs.pakid").
		Joins("join mk on mk.mkid = krs.mkid").
		Joins("join jad on jad.id_master_kegiatan = vw_kelas_tawar.id_master_kegiatan and mk.jurid = jad.jurid").
		Joins("join pak on pak.pakid = krs.pakid").
		Where("krs.mhsid = ?", mhsID).
		Where("pak.pakid = ?", pakID).
		Order("mk.mknama ASC").
		Find(&studentSubject).Error

	return studentSubject, err
}

func (s *studentSubjectRepository) GetSubjectClassStudentWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectSchedule, error) {
	subjectSchedules := []entities.SubjectSchedule{}

	err := s.db.WithContext(ctx).Table("krs").
		Distinct("krs.mkid", "jad.kultipeid", "jad.ruangid", "hari.haridesc", "jad.hari as hari_code", "jam.jammulai", "jam.jamhingga", "jad.sks").
		Joins("join mk on mk.mkid = krs.mkid").
		Joins("join pak on pak.pakid = krs.pakid").
		Joins("join jad on mk.mkid = jad.mkid and pak.pakid = jad.pakid and krs.kelaskrs = jad.kelas").
		Joins("join jam on jam.jamid = jad.jamid").
		Joins("join hari on hari.hari = jad.hari").
		Where("krs.mhsid = ?", mhsID).
		Where("pak.pakid = ?", pakID).
		Find(&subjectSchedules).Error

	return subjectSchedules, err
}

func (s *studentSubjectRepository) StudyPeriodes(ctx context.Context, mhsID string) ([]entities.AcademicPeriod, error) {
	academicPeriods := []entities.AcademicPeriod{}

	err := s.db.WithContext(ctx).
		Table("pak").
		Select("pak.*").
		Joins("JOIN krs on pak.pakid = krs.pakid").
		Where("krs.mhsid = ?", mhsID).
		Group("pak.pakid").
		Order("pak.pakid DESC").
		Find(&academicPeriods).Error

	return academicPeriods, err
}

func (s *studentSubjectRepository) GetActiveRoomLoans(ctx context.Context) ([]entities.RoomKeyLoan, error) {
	loans := []entities.RoomKeyLoan{}

	err := s.db.WithContext(ctx).
		Table("peminjaman_ruang").
		Where("waktu_kembali IS NULL").
		Order("waktu_pinjam DESC, id_peminjaman_ruang DESC").
		Find(&loans).Error

	return loans, err
}

func (s *studentSubjectRepository) GetTodayHariCode(ctx context.Context, dayOfWeek int) (string, error) {
	var result entities.TodayHariCode

	err := s.db.WithContext(ctx).
		Table("hari").
		Select("hari").
		Where("day = ?", strconv.Itoa(dayOfWeek)).
		First(&result).Error

	if err != nil {
		return "", err
	}

	return result.HariCode, nil
}
