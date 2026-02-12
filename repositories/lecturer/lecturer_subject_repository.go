package lecturer_repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type lecturerSubjectRepository struct {
	db gorm.DB
}

type LecturerSubjectRepository interface {
	GetActiveSubjectByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubject, error)
	LecturePeriodes(ctx context.Context, dosID string) ([]entities.AcademicPeriod, error)
	GetSubjectByLecturerIDWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.LecturerSubject, error)
	GetLecturerSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubject, error)
	LecturerSubjectMajor(ctx context.Context, dosID string, pakID string) ([]entities.Major, error)
	GetSubjectClassLecturer(ctx context.Context, dosID string) ([]entities.SubjectSchedule, error)
	GetSubjectClassLecturerWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.SubjectSchedule, error)
	GetSubjectClassLecturerFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectSchedule, error)
	GetStudentScore(ctx context.Context, pakid string, mkid string, class string) ([]entities.StudentScore, error)
	GetAssignment(ctx context.Context, masterActivityID string) ([]entities.Assignment, error)
	GetAssignmentSubmission(ctx context.Context, masterActivityID string) ([]entities.AssignmentSubmission, error)
	GetClassOffered(ctx context.Context, masterActivityID string) (entities.ClassOffered, error)
	GetActiveSubjectReportByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubjectReport, error)
	GetSubjectReportByLecturerIDWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.LecturerSubjectReport, error)
	GetLecturerReportSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubjectReport, error)
}

func NewLecturerSubjectRepository(db *gorm.DB) *lecturerSubjectRepository {
	return &lecturerSubjectRepository{db: *db}
}

func (l *lecturerSubjectRepository) GetActiveSubjectByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubject, error) {
	lecturerSubjects := []entities.LecturerSubject{}
	var pakID string

	err := l.db.WithContext(ctx).Table("(?) as pak", l.db.Table("pak").Order("pakid DESC").Limit(5)).Select("pakid").Where("isactive = ?", true).Find(&pakID).Error

	if err != nil {
		return lecturerSubjects, err
	}

	err = l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?)
            ORDER BY mknama, kelas;`, pakID, dosID, dosID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}

func (l *lecturerSubjectRepository) GetSubjectClassLecturer(ctx context.Context, dosID string) ([]entities.SubjectSchedule, error) {
	subjectSchedules := []entities.SubjectSchedule{}
	var pakID string

	err := l.db.WithContext(ctx).Table("(?) as pak", l.db.Table("pak").Order("pakid DESC").Limit(5)).Select("pakid").Where("isactive = ?", true).Find(&pakID).Error

	if err != nil {
		return subjectSchedules, err
	}

	err = l.db.WithContext(ctx).Table("vw_jadwal").
		Distinct("vw_jadwal.mkid", "vw_jadwal.kelas", "vw_jadwal.kultipeid", "vw_jadwal.ruangid", "vw_jadwal.haridesc", "vw_jadwal.jammulai", "vw_jadwal.jamhingga").
		Joins("join jad on jad.mkid = vw_jadwal.mkid and jad.pakid = vw_jadwal.pakid and vw_jadwal.kelas = jad.kelas").
		Where("jad.dosid = ?", dosID).
		Where("jad.pakid = ?", pakID).
		Find(&subjectSchedules).Error

	return subjectSchedules, err
}

func (l *lecturerSubjectRepository) GetSubjectByLecturerIDWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.LecturerSubject, error) {
	lecturerSubjects := []entities.LecturerSubject{}

	err := l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?)
            ORDER BY mknama, kelas;`, pakID, dosID, dosID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}

func (l *lecturerSubjectRepository) GetSubjectClassLecturerWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.SubjectSchedule, error) {
	subjectSchedules := []entities.SubjectSchedule{}

	err := l.db.WithContext(ctx).Table("vw_jadwal").
		Distinct("vw_jadwal.mkid", "vw_jadwal.kelas", "vw_jadwal.kultipeid", "vw_jadwal.ruangid", "vw_jadwal.haridesc", "vw_jadwal.jammulai", "vw_jadwal.jamhingga").
		Joins("join jad on jad.mkid = vw_jadwal.mkid and jad.pakid = vw_jadwal.pakid and vw_jadwal.kelas = jad.kelas").
		Where("jad.dosid = ?", dosID).
		Where("jad.pakid = ?", pakID).
		Find(&subjectSchedules).Error

	return subjectSchedules, err
}

func (l *lecturerSubjectRepository) GetLecturerSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubject, error) {
	lecturerSubjects := []entities.LecturerSubject{}

	err := l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?) AND jad.jurasli_jurid = ?
            ORDER BY mknama, kelas;`, pakID, dosID, dosID, jurID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}

func (l *lecturerSubjectRepository) GetSubjectClassLecturerFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectSchedule, error) {
	subjectSchedules := []entities.SubjectSchedule{}

	err := l.db.WithContext(ctx).Table("vw_jadwal").
		Distinct("vw_jadwal.mkid", "vw_jadwal.kelas", "vw_jadwal.kultipeid", "vw_jadwal.ruangid", "vw_jadwal.haridesc", "vw_jadwal.jammulai", "vw_jadwal.jamhingga").
		Joins("join jad on jad.mkid = vw_jadwal.mkid and jad.pakid = vw_jadwal.pakid and vw_jadwal.kelas = jad.kelas").
		Where("jad.dosid = ?", dosID).
		Where("jad.pakid = ?", pakID).
		Where("vw_jadwal.jurasli_jurid = ?", jurID).
		Find(&subjectSchedules).Error

	return subjectSchedules, err
}

func (l *lecturerSubjectRepository) LecturePeriodes(ctx context.Context, dosID string) ([]entities.AcademicPeriod, error) {
	academicPeriods := []entities.AcademicPeriod{}

	err := l.db.WithContext(ctx).
		Table("jad").
		Distinct("pak.*").
		Joins("JOIN pak on pak.pakid = jad.pakid").
		Where("jad.dosid = ?", dosID).
		Order("pak.pakid DESC").
		Limit(5).
		Find(&academicPeriods).Error

	return academicPeriods, err
}

func (l *lecturerSubjectRepository) LecturerSubjectMajor(ctx context.Context, dosID string, pakID string) ([]entities.Major, error) {
	lectureMajor := []entities.Major{}

	err := l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jur.jurid, jad.jurasli_jurid, jur.jurnama, jur.prodinama
			FROM vw_jadwal as jad
			JOIN jur ON jur.jurid = jad.jurasli_jurid
			JOIN (SELECT * FROM pak ORDER BY pakid DESC LIMIT 5) AS pak ON pak.pakid = jad.pakid
			WHERE jad.dosid = ? AND pak.pakid = ?;`, dosID, pakID).
		// Raw(`SELECT DISTINCT jur.jurid, jad.jurasli_jurid, jur.jurnama, jur.prodinama
		// 	FROM vw_jadwal as jad
		// 	JOIN jur ON jur.jurid = jad.jurasli_jurid
		// 	JOIN (SELECT * FROM pak ORDER BY pakid DESC LIMIT 5) AS pak ON pak.pakid = jad.pakid
		// 	WHERE jad.dosid = ?`, dosID).
		Find(&lectureMajor).Error

	return lectureMajor, err
}

func (l *lecturerSubjectRepository) GetStudentScore(ctx context.Context, pakid string, mkid string, class string) ([]entities.StudentScore, error) {
	studentScore := []entities.StudentScore{}

	err := l.db.WithContext(ctx).Table("krs").Select("krs.mhsid", "krs.nilagk", "krs.nilhrf", "mhs.mhsnama").
		Joins("join mhs on mhs.mhsid = krs.mhsid").
		Where("krs.pakid = ?", pakid).
		Where("krs.mkid = ?", mkid).
		Where("krs.kelaskrs = ?", class).
		Order("krs.mhsid DESC").Find(&studentScore).Error

	return studentScore, err
}

func (l *lecturerAssignmentRepository) GetReg(ctx context.Context) {
	// err = l.db.WithContext(ctx).Table("reg").Where("regid = ?", "NIL_TNGL_START").Error
	// err = l.db.WithContext(ctx).Table("reg").Where("regid = ?", "NIL_TNGL_END").Error
}

func (l *lecturerSubjectRepository) GetAssignment(ctx context.Context, masterActivityID string) ([]entities.Assignment, error) {
	assignments := []entities.Assignment{}

	err := l.db.WithContext(ctx).Table("tugas_kul").Select("*").
		Joins("join jnil on tugas_kul.jnilid = jnil.jnilid").
		Where("master_kegiatan_id = ?", masterActivityID).
		Find(&assignments).Error

	return assignments, err
}

func (l *lecturerSubjectRepository) GetAssignmentSubmission(ctx context.Context, masterActivityID string) ([]entities.AssignmentSubmission, error) {
	studentAssignmentSubmissions := []entities.AssignmentSubmission{}

	err := l.db.WithContext(ctx).Table("tugas_submission").
		Joins("join tugas_kul on tugas_kul.id_tugas_kul = tugas_submission.tugas_kul_id").
		Where("master_kegiatan_id = ?", masterActivityID).
		Find(&studentAssignmentSubmissions).Error

	return studentAssignmentSubmissions, err
}

func (l *lecturerSubjectRepository) GetClassOffered(ctx context.Context, masterActivityID string) (entities.ClassOffered, error) {
	ClassOffered := entities.ClassOffered{}

	err := l.db.WithContext(ctx).Table("klstw").
		Where("id_master_kegiatan = ?", masterActivityID).
		Find(&ClassOffered).Error

	return ClassOffered, err
}

func (l *lecturerSubjectRepository) GetActiveSubjectReportByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubjectReport, error) {
	lecturerSubjects := []entities.LecturerSubjectReport{}
	var pakID string

	err := l.db.WithContext(ctx).Table("(?) as pak", l.db.Table("pak").Order("pakid DESC").Limit(5)).Select("pakid").Where("isactive = ?", true).Find(&pakID).Error

	if err != nil {
		return lecturerSubjects, err
	}

	err = l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan, hari.haridesc, jam.jammulai, jam.jamhingga, jad.kultipeid, jam.jamid, jad.ruangid
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?)
            ORDER BY mknama, kelas;`, pakID, dosID, dosID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}

func (l *lecturerSubjectRepository) GetSubjectReportByLecturerIDWithPeriod(ctx context.Context, dosID string, pakID string) ([]entities.LecturerSubjectReport, error) {
	lecturerSubjects := []entities.LecturerSubjectReport{}

	err := l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan,  hari.haridesc, jam.jammulai, jam.jamhingga, jad.kultipeid, jam.jamid, jad.ruangid
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?)
            ORDER BY mknama, kelas;`, pakID, dosID, dosID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}

func (l *lecturerSubjectRepository) GetLecturerReportSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubjectReport, error) {
	lecturerSubjects := []entities.LecturerSubjectReport{}

	err := l.db.WithContext(ctx).
		Raw(`SELECT DISTINCT jad.pakid, jad.dosid, jad.dosnama, mkid, mknama, kelas, mksks, jurasli_jurid AS jurid, jurasli_jurnama, (SELECT COUNT(*) FROM krs WHERE pakid = jad.pakid AND mkid = jad.mkid AND kelaskrs = jad.kelas) AS jumlah_mhs_perkelas, id_master_kegiatan,  hari.haridesc, jam.jammulai, jam.jamhingga, jad.kultipeid, jam.jamid, jad.ruangid
            FROM vw_jadwal AS jad
            LEFT JOIN jadteam ON jad.id_master_kegiatan = jadteam.id_jad_master
			join jam on jam.jamid = jad.jamid
			join hari on hari.hari = jad.hari
            WHERE jad.pakid = ? AND (jad.dosid = ? OR jadteam.dosid = ?) AND jad.jurasli_jurid = ?
            ORDER BY mknama, kelas;`, pakID, dosID, dosID, jurID).
		Find(&lecturerSubjects).Error

	return lecturerSubjects, err
}
