package student_repositories

import (
	"classroom_itats_api/entities"
	"context"
	"time"

	"gorm.io/gorm"
)

type studentMaterialRepository struct {
	db *gorm.DB
}

type StudentMaterialRepository interface {
	GetWeekMaterial(ctx context.Context, pakID string, mkID string, class string) ([]entities.LectureWeek, error)
	GetStudyAchievement(ctx context.Context, pakID string, mkID string, class string) ([]entities.StudyAchievement, error)
	GetStudentMaterial(ctx context.Context, mkID string, class string, pakID string, weekID int) ([]entities.Material, error)
	GetStudentAssignment(ctx context.Context, masterActivityID string, weekID float64, mhsID string) ([]entities.Assignment, error)
	GetStudentAssignmentGroup(ctx context.Context, masterActivityID string) ([]entities.Assignment, error)
	GetStudentScore(ctx context.Context, mhsID string, masterActivityID string) ([]entities.StudentAssignmentScore, error)
	GetStudentAssignmentSubmission(ctx context.Context, mhsID string, assignmentID int) (entities.AssignmentSubmission, error)
	AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error)
	AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error)
	GetActiveAssignment(ctx context.Context, masterActivityID string, mhsID string) ([]entities.Assignment, error)
	GetHomeActiveAssignment(ctx context.Context, mhsID string, pakID string) ([]entities.Assignment, error)
}

func NewStudentMaterialRepository(db *gorm.DB) *studentMaterialRepository {
	return &studentMaterialRepository{db: db}
}

func (s *studentMaterialRepository) GetWeekMaterial(ctx context.Context, pakID string, mkID string, class string) ([]entities.LectureWeek, error) {
	LectureWeek := []entities.LectureWeek{}

	err := s.db.WithContext(ctx).Raw(`
		select kul.mkid, kul.kelas, kul.pakid, kul.jurid, kul.dosid, kul.weekid, kul.kulid, kul.link_meet, kul.jenis_kuliah, kul.link_record, jenis_perkuliahan.nama_jenis_perkuliahan from kul
		join jenis_perkuliahan on jenis_perkuliahan.id_jenis_perkuliahan = kul.jenis_kuliah
        where kul.mkid = ?
		and kul.kelas = ?
		and kul.pakid = ?
		order by kul.weekid asc
	`, mkID, class, pakID).Find(&LectureWeek).Error

	return LectureWeek, err
}

func (s *studentMaterialRepository) GetStudyAchievement(ctx context.Context, pakID string, mkID string, class string) ([]entities.StudyAchievement, error) {
	studyAchievements := []entities.StudyAchievement{}

	err := s.db.WithContext(ctx).Raw(`
		select capaian_pembelajaran.*, rencana_pembelajaran.*, mapping_rps.* from capaian_pembelajaran
		join rencana_pembelajaran on rencana_pembelajaran.capaian_pembelajaran_id = capaian_pembelajaran.id_capaian_pembelajaran
		join mapping_rps on mapping_rps.id_rencana_pembelajaran = rencana_pembelajaran.id_rencana_pembelajaran
		join kul on kul.dosid = rencana_pembelajaran.dosid and kul.jurid = rencana_pembelajaran.jurid and
		kul.kelas = rencana_pembelajaran.kelas and kul.mkid = rencana_pembelajaran.mkid and
		kul.pakid = rencana_pembelajaran.pakid and kul.weekid = mapping_rps.weekid
		where capaian_pembelajaran.mkid = ?
		and rencana_pembelajaran.kelas = ?
		and rencana_pembelajaran.pakid = ?
		order by mapping_rps.weekid asc
	`, mkID, class, pakID).Find(&studyAchievements).Error

	return studyAchievements, err
}

func (s *studentMaterialRepository) GetStudentMaterial(ctx context.Context, mkID string, class string, pakID string, weekID int) ([]entities.Material, error) {
	studentMaterials := []entities.Material{}

	err := s.db.WithContext(ctx).Raw(`
		select * from materi
		join kul_materi on kul_materi.materi_id = materi.materi_id
		join kul on kul.kulid = kul_materi.kul_id
		where kul.mkid = ? 
		and kul.kelas = ?
		and kul.pakid = ?
		and kul.weekid = ?
		and hidden_status = ?
		and deleted_at is null
	`, mkID, class, pakID, weekID, 0).Find(&studentMaterials).Error

	return studentMaterials, err
}

func (s *studentMaterialRepository) GetStudentAssignmentGroup(ctx context.Context, masterActivityID string) ([]entities.Assignment, error) {
	studentAssignments := []entities.Assignment{}

	err := s.db.WithContext(ctx).Table("tugas_kul").Select("tugas_kul.master_kegiatan_id, tugas_kul.weekid").
		Joins("join jnil on tugas_kul.jnilid = jnil.jnilid").
		Where("master_kegiatan_id = ?", masterActivityID).Group("master_kegiatan_id, weekid").
		Find(&studentAssignments).Error

	return studentAssignments, err
}

func (s *studentMaterialRepository) GetStudentAssignment(ctx context.Context, masterActivityID string, weekID float64, mhsID string) ([]entities.Assignment, error) {
	studentAssignments := []entities.Assignment{}

	err := s.db.WithContext(ctx).Table("tugas_kul").
		Select("tugas_kul.*").
		Joins("join jnil on tugas_kul.jnilid = jnil.jnilid").
		Where("master_kegiatan_id = ?", masterActivityID).Where("weekid = ?", weekID).
		Find(&studentAssignments).Error

	return studentAssignments, err
}

func (s *studentMaterialRepository) GetStudentAssignmentSubmission(ctx context.Context, mhsID string, assignmentID int) (entities.AssignmentSubmission, error) {
	studentAssignmentSubmissions := entities.AssignmentSubmission{}

	err := s.db.WithContext(ctx).Table("tugas_submission").
		Where("mhsid = ?", mhsID).
		Where("tugas_kul_id = ?", assignmentID).
		Find(&studentAssignmentSubmissions).Error

	return studentAssignmentSubmissions, err
}

func (s *studentMaterialRepository) GetStudentScore(ctx context.Context, mhsID string, masterActivityID string) ([]entities.StudentAssignmentScore, error) {
	studentAssignmentScores := []entities.StudentAssignmentScore{}

	err := s.db.WithContext(ctx).Table("tugas_submission").
		Select("tugas_submission.*, tugas_kul.judul_tugas").
		Joins("join tugas_kul on tugas_kul.id_tugas_kul = tugas_submission.tugas_kul_id").
		Where("mhsid = ?", mhsID).
		Where("master_kegiatan_id = ?", masterActivityID).
		Where("tugas_kul.is_tampil = ?", true).
		Find(&studentAssignmentScores).Error

	return studentAssignmentScores, err
}

func (s *studentMaterialRepository) AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	tugasKul := []entities.Assignment{}

	err := s.db.WithContext(ctx).Table("tugas_kul").Where("created_at between ? and ?", time.Now().Add(time.Duration(-30)*time.Minute), time.Now()).Find(&tugasKul).Error
	if err != nil {
		return nil, err
	}

	for _, tgskul := range tugasKul {
		classOffered := entities.ClassOffered{}
		var subject string
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string
		e := s.db.WithContext(ctx).Table("vw_kelas_tawar").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error

		if e != nil {
			err = e
		}

		e = s.db.WithContext(ctx).Table("jur").Select("jurid").Where("jurid = ?", classOffered.MajorID).Or("jur_parent_id = ?", classOffered.MajorID).Find(&jurID).Error

		if e != nil {
			err = e
		}

		if len(jurID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", classOffered.AcademicPeriodID).
				Where("mkid = ?", classOffered.SubjectID).
				Where("kelaskrs = ?", classOffered.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT mhsid FROM tugas_submission where tugas_submission.tugas_kul_id  = ?)", tgskul.AssignmentID).
				Find(&mhsID).Error
		}

		if e != nil {
			err = e
		}

		if len(mhsID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("users").Select("mobile_token").Where("name in ?", mhsID).Where("mobile_token != ?", "null").Find(&user).Error
		}

		if e == nil {
			usr["tugaskul"] = tgskul
			usr["klstw"] = classOffered
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

func (s *studentMaterialRepository) AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	tugasKul := []entities.Assignment{}

	err := s.db.WithContext(ctx).Table("tugas_kul").Where("batas_pengumpulan between ? and ?", time.Now(), time.Now().Add(time.Duration(3)*time.Hour)).Find(&tugasKul).Error
	if err != nil {
		return nil, err
	}

	for _, tgskul := range tugasKul {
		classOffered := entities.ClassOffered{}
		var subject string
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string
		e := s.db.WithContext(ctx).Table("vw_kelas_tawar").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).Find(&classOffered).Error

		if e != nil {
			err = e
		}

		e = s.db.WithContext(ctx).Table("jur").Select("jurid").Where("jurid = ?", classOffered.MajorID).Or("jur_parent_id = ?", classOffered.MajorID).Find(&jurID).Error

		if e != nil {
			err = e
		}

		if len(jurID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", classOffered.AcademicPeriodID).
				Where("mkid = ?", classOffered.SubjectID).
				Where("kelaskrs = ?", classOffered.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT mhsid FROM tugas_submission where tugas_submission.tugas_kul_id  = ?)", tgskul.AssignmentID).Find(&mhsID).Error
		}

		if e != nil {
			err = e
		}

		if len(mhsID) > 0 || e == nil {
			e = s.db.WithContext(ctx).Table("users").Select("mobile_token").Where("name in ?", mhsID).Where("mobile_token != ?", "null").Find(&user).Error
		}

		if e == nil {
			usr["tugaskul"] = tgskul
			usr["klstw"] = classOffered
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

func (s *studentMaterialRepository) GetActiveAssignment(ctx context.Context, masterActivityID string, mhsID string) ([]entities.Assignment, error) {
	activeAssignments := []entities.Assignment{}
	today := time.Now()

	err := s.db.WithContext(ctx).Table("tugas_kul").
		Select("tugas_kul.*").
		Joins("join jnil on tugas_kul.jnilid = jnil.jnilid").
		Where("tugas_kul.master_kegiatan_id = ?", masterActivityID).
		Where("tugas_kul.waktu_mulai_tugas <= ?", today).
		Where("tugas_kul.waktu_akhir_tugas >= ?", today).
		Where("NOT EXISTS (SELECT 1 FROM tugas_submission WHERE tugas_submission.tugas_kul_id = tugas_kul.id_tugas_kul AND tugas_submission.mhsid = ?)", mhsID).
		Find(&activeAssignments).Error

	return activeAssignments, err
}

func (s *studentMaterialRepository) GetHomeActiveAssignment(ctx context.Context, mhsID string, pakID string) ([]entities.Assignment, error) {
	activeAssignments := []entities.Assignment{}

	// Get active pakid first if not provided
	if pakID == "" {
		err := s.db.WithContext(ctx).Table("(?) as pak", s.db.Table("pak").Order("pakid DESC").Limit(5)).
			Select("pakid").Where("isactive = ?", true).Find(&pakID).Error
		if err != nil {
			return nil, err
		}
	}

	// UNION query: support both old periods (via vw_kelas_tawar) and new periods (via jad)
	rawSQL := `
		SELECT DISTINCT ON (tugas_kul.id_tugas_kul)
			tugas_kul.*,
			mk.mknama,
			src.kelas
		FROM tugas_kul
		JOIN (
			-- Jalur lama: via vw_kelas_tawar
			SELECT vw.id_master_kegiatan, vw.kelas, krs.mhsid, krs.pakid
			FROM vw_kelas_tawar vw
			JOIN krs ON krs.mkid = vw.mkid AND krs.kelaskrs = vw.kelas AND krs.pakid = vw.pakid
			UNION
			-- Jalur baru: via jad
			SELECT jad.id_master_kegiatan, jad.kelas, krs.mhsid, krs.pakid
			FROM jad
			JOIN krs ON krs.mkid = jad.mkid AND krs.kelaskrs = jad.kelas AND krs.pakid = jad.pakid
		) src ON src.id_master_kegiatan = tugas_kul.master_kegiatan_id
		LEFT JOIN mk ON mk.mkid = COALESCE(
			(SELECT mkid FROM vw_kelas_tawar WHERE id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1),
			(SELECT mkid FROM jad WHERE id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1)
		)
		WHERE src.mhsid = ?
		  AND src.pakid = ?
		  AND (tugas_kul.waktu_mulai_tugas IS NULL OR tugas_kul.waktu_mulai_tugas <= NOW())
		  AND (tugas_kul.waktu_akhir_tugas IS NULL OR tugas_kul.waktu_akhir_tugas >= NOW())
		  AND NOT EXISTS (
				SELECT 1 FROM tugas_submission
				WHERE tugas_submission.tugas_kul_id = tugas_kul.id_tugas_kul
				  AND tugas_submission.mhsid = ?
		  )
	`

	err := s.db.WithContext(ctx).Raw(rawSQL, mhsID, pakID, mhsID).Scan(&activeAssignments).Error

	return activeAssignments, err
}
