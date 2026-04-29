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
	GetStudentAssignmentGroup(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Assignment, error)
	GetStudentScore(ctx context.Context, mhsID string, pakID string, mkID string, class string) ([]entities.StudentAssignmentScore, error)
	GetStudentAssignmentSubmission(ctx context.Context, mhsID string, assignmentID int) (entities.AssignmentSubmission, error)
	AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error)
	AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error)
	AssignmentReminderH1(ctx context.Context) ([]map[string]interface{}, error)
	GetActiveAssignment(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Assignment, error)
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

func (s *studentMaterialRepository) GetStudentAssignmentGroup(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Assignment, error) {
	studentAssignments := []entities.Assignment{}

	// Meniru logika PHP getTugasByMatkul:
	// Cari master_kegiatan_id dari 2 jalur berdasarkan KRS mahasiswa
	rawSQL := `
		SELECT DISTINCT ON (tugas_kul.weekid)
			tugas_kul.*,
			COALESCE(
				(SELECT jad.kelas FROM jad WHERE jad.id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1),
				(SELECT vw.kelas FROM vw_kelas_tawar vw WHERE vw.id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1)
			) AS kelas,
			mk.mknama,
			(ts.id_tugas_submission IS NOT NULL) AS sudah_submit,
			ts.file_tugas AS submission_file,
			ts.link_tugas AS submission_link,
			ts.created_at AS submission_date
		FROM tugas_kul
		LEFT JOIN mk ON mk.mkid = COALESCE(
			(SELECT jad2.mkid FROM jad jad2 WHERE jad2.id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1),
			(SELECT vw2.mkid FROM vw_kelas_tawar vw2 WHERE vw2.id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1)
		)
		LEFT JOIN tugas_submission ts ON ts.tugas_kul_id = tugas_kul.id_tugas_kul AND ts.mhsid = ?
		WHERE tugas_kul.master_kegiatan_id IN (
			-- Jalur lama: via vw_kelas_tawar + krs
			SELECT vw.id_master_kegiatan
			FROM vw_kelas_tawar vw
			JOIN krs ON krs.mkid = vw.mkid AND krs.kelaskrs = vw.kelas AND krs.pakid = vw.pakid
			WHERE krs.mhsid = ?
			  AND krs.pakid = ?
			  AND krs.mkid = ?

			UNION

			-- Jalur baru (20252+): via jad + krs
			SELECT jad.id_master_kegiatan
			FROM jad
			JOIN krs ON krs.mkid = jad.mkid AND krs.kelaskrs = jad.kelas AND krs.pakid = jad.pakid
			WHERE krs.mhsid = ?
			  AND krs.pakid = ?
			  AND krs.mkid = ?
		)
		ORDER BY tugas_kul.weekid, tugas_kul.id_tugas_kul
	`

	err := s.db.WithContext(ctx).Raw(rawSQL,
		mhsID,       // LEFT JOIN tugas_submission ts.mhsid
		mhsID, pakID, mkID, // jalur lama: krs.mhsid, krs.pakid, krs.mkid
		mhsID, pakID, mkID, // jalur baru: krs.mhsid, krs.pakid, krs.mkid
	).Scan(&studentAssignments).Error

	return studentAssignments, err
}

func (s *studentMaterialRepository) GetStudentAssignment(ctx context.Context, masterActivityID string, weekID float64, mhsID string) ([]entities.Assignment, error) {
	studentAssignments := []entities.Assignment{}

	err := s.db.WithContext(ctx).Table("tugas_kul").
		Select("tugas_kul.*, (ts.id_tugas_submission IS NOT NULL) AS sudah_submit, ts.file_tugas AS submission_file, ts.link_tugas AS submission_link, ts.created_at AS submission_date").
		Joins("join jnil on tugas_kul.jnilid = jnil.jnilid").
		Joins("LEFT JOIN tugas_submission ts ON ts.tugas_kul_id = tugas_kul.id_tugas_kul AND ts.mhsid = ?", mhsID).
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

func (s *studentMaterialRepository) GetStudentScore(ctx context.Context, mhsID string, pakID string, mkID string, class string) ([]entities.StudentAssignmentScore, error) {
	studentAssignmentScores := []entities.StudentAssignmentScore{}

	rawSQL := `
		SELECT tugas_submission.*, tugas_kul.judul_tugas
		FROM tugas_submission
		JOIN tugas_kul ON tugas_kul.id_tugas_kul = tugas_submission.tugas_kul_id
		JOIN (
			-- Semester lama: via vw_kelas_tawar
			SELECT id_master_kegiatan FROM vw_kelas_tawar
			WHERE mkid = ? AND kelas = ? AND pakid = ?
			UNION
			-- Semester baru (20252+): via jad
			SELECT id_master_kegiatan FROM jad
			WHERE mkid = ? AND kelas = ? AND pakid = ?
		) src ON src.id_master_kegiatan = tugas_kul.master_kegiatan_id
		WHERE tugas_submission.mhsid = ?
		  AND tugas_kul.is_tampil = true
	`

	err := s.db.WithContext(ctx).Raw(rawSQL, mkID, class, pakID, mkID, class, pakID, mhsID).Scan(&studentAssignmentScores).Error

	return studentAssignmentScores, err
}

func (s *studentMaterialRepository) AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	tugasKul := []entities.Assignment{}

	// Ambil tugas yang dibuat dalam 60 menit terakhir (lebih longgar dari 30 menit untuk antisipasi delay)
	err := s.db.WithContext(ctx).Table("tugas_kul").
		Where("created_at between ? and ?", time.Now().Add(time.Duration(-60)*time.Minute), time.Now()).
		Order("created_at DESC").
		Find(&tugasKul).Error
	if err != nil {
		return nil, err
	}

	for _, tgskul := range tugasKul {
		classOffered := entities.ClassOffered{}
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string

		// Coba cari di vw_kelas_tawar
		e := s.db.WithContext(ctx).Table("vw_kelas_tawar").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		
		// Jika tidak ada di vw_kelas_tawar, coba di tabel jad (untuk semester 20252+)
		if e != nil {
			e = s.db.WithContext(ctx).Table("jad").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		}

		if e != nil {
			// Jika tetap tidak ditemukan, lewati tugas ini (mungkin data master_kegiatan tidak valid)
			continue
		}

		// Cari Jurusan (termasuk parent/child)
		s.db.WithContext(ctx).Table("jur").Select("jurid").
			Where("jurid = ?", classOffered.MajorID).
			Or("jur_parent_id = ?", classOffered.MajorID).
			Find(&jurID)

		if len(jurID) > 0 {
			// Cari Mahasiswa yang mengambil matakuliah tersebut di KRS dan belum submit tugas ini
			s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", classOffered.AcademicPeriodID).
				Where("mkid = ?", classOffered.SubjectID).
				Where("kelaskrs = ?", classOffered.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT 1 FROM tugas_submission ts WHERE ts.tugas_kul_id = ? AND ts.mhsid = krs.mhsid)", tgskul.AssignmentID).
				Find(&mhsID)
		}

		if len(mhsID) > 0 {
			var mhsTokens []struct {
				Name        string
				MobileToken string
			}
			s.db.WithContext(ctx).Table("users").Select("name, mobile_token").
				Where("name in ?", mhsID).
				Where("mobile_token IS NOT NULL AND mobile_token != ? AND mobile_token != ?", "null", "").
				Find(&mhsTokens)

			if len(mhsTokens) > 0 {
				tokenMap := make(map[string]string)
				var tokens []string
				for _, mt := range mhsTokens {
					tokenMap[mt.Name] = mt.MobileToken
					tokens = append(tokens, mt.MobileToken)
				}

				usr := map[string]interface{}{}
				usr["tugaskul"] = tgskul
				usr["klstw"] = classOffered
				usr["user"] = tokens
				usr["token_map"] = tokenMap
				usr["subject"] = classOffered.SubjectName
				users = append(users, usr)
			}
		}
	}

	return users, nil
}

func (s *studentMaterialRepository) AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	tugasKul := []entities.Assignment{}

	// Cari tugas yang deadline-nya antara sekarang sampai 3 jam ke depan
	err := s.db.WithContext(ctx).Table("tugas_kul").
		Where("batas_pengumpulan between ? and ?", time.Now(), time.Now().Add(time.Duration(3)*time.Hour)).
		Find(&tugasKul).Error
	if err != nil {
		return nil, err
	}

	for _, tgskul := range tugasKul {
		classOffered := entities.ClassOffered{}
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string

		// Coba vw_kelas_tawar
		e := s.db.WithContext(ctx).Table("vw_kelas_tawar").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		
		// Fallback ke jad
		if e != nil {
			e = s.db.WithContext(ctx).Table("jad").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		}

		if e != nil {
			continue
		}

		s.db.WithContext(ctx).Table("jur").Select("jurid").
			Where("jurid = ?", classOffered.MajorID).
			Or("jur_parent_id = ?", classOffered.MajorID).
			Find(&jurID)

		if len(jurID) > 0 {
			s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", classOffered.AcademicPeriodID).
				Where("mkid = ?", classOffered.SubjectID).
				Where("kelaskrs = ?", classOffered.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT 1 FROM tugas_submission ts WHERE ts.tugas_kul_id = ? AND ts.mhsid = krs.mhsid)", tgskul.AssignmentID).
				Find(&mhsID)
		}

		if len(mhsID) > 0 {
			s.db.WithContext(ctx).Table("users").
				Select("mobile_token").
				Where("name in ?", mhsID).
				Where("mobile_token IS NOT NULL AND mobile_token != ? AND mobile_token != ?", "null", "").
				Find(&user)
		}

		if len(user) > 0 {
			usr["tugaskul"] = tgskul
			usr["klstw"] = classOffered
			usr["user"] = user
			usr["subject"] = classOffered.SubjectName
			users = append(users, usr)
		}
	}

	return users, nil
}

func (s *studentMaterialRepository) AssignmentReminderH1(ctx context.Context) ([]map[string]interface{}, error) {
	users := []map[string]interface{}{}
	tugasKul := []entities.Assignment{}

	// Cari tugas yang deadline-nya besok (H-1)
	err := s.db.WithContext(ctx).Table("tugas_kul").
		Where("DATE(batas_pengumpulan) = CURRENT_DATE + INTERVAL '1 day'").
		Find(&tugasKul).Error
	if err != nil {
		return nil, err
	}

	for _, tgskul := range tugasKul {
		classOffered := entities.ClassOffered{}
		var user []string
		usr := map[string]interface{}{}
		var mhsID []string
		var jurID []string

		// Coba vw_kelas_tawar
		e := s.db.WithContext(ctx).Table("vw_kelas_tawar").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		
		// Fallback ke jad
		if e != nil {
			e = s.db.WithContext(ctx).Table("jad").Where("id_master_kegiatan = ?", tgskul.ActivityMasterID).First(&classOffered).Error
		}

		if e != nil {
			continue
		}

		s.db.WithContext(ctx).Table("jur").Select("jurid").
			Where("jurid = ?", classOffered.MajorID).
			Or("jur_parent_id = ?", classOffered.MajorID).
			Find(&jurID)

		if len(jurID) > 0 {
			s.db.WithContext(ctx).Table("krs").
				Select("mhsid").
				Where("pakid = ?", classOffered.AcademicPeriodID).
				Where("mkid = ?", classOffered.SubjectID).
				Where("kelaskrs = ?", classOffered.SubjectClass).
				Where("jurid in ?", jurID).
				Where("NOT EXISTS (SELECT 1 FROM tugas_submission ts WHERE ts.tugas_kul_id = ? AND ts.mhsid = krs.mhsid)", tgskul.AssignmentID).
				Find(&mhsID)
		}

		if len(mhsID) > 0 {
			s.db.WithContext(ctx).Table("users").
				Select("mobile_token").
				Where("name in ?", mhsID).
				Where("mobile_token IS NOT NULL AND mobile_token != ? AND mobile_token != ?", "null", "").
				Find(&user)
		}

		if len(user) > 0 {
			usr["tugaskul"] = tgskul
			usr["klstw"] = classOffered
			usr["user"] = user
			usr["subject"] = classOffered.SubjectName
			users = append(users, usr)
		}
	}

	return users, nil
}

func (s *studentMaterialRepository) GetActiveAssignment(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Assignment, error) {
	activeAssignments := []entities.Assignment{}

	rawSQL := `
		SELECT DISTINCT ON (tugas_kul.id_tugas_kul)
			tugas_kul.*,
			src.kelas,
			mk.mknama,
			(ts.id_tugas_submission IS NOT NULL) AS sudah_submit,
			ts.file_tugas AS submission_file,
			ts.link_tugas AS submission_link,
			ts.created_at AS submission_date
		FROM tugas_kul
		JOIN jnil ON tugas_kul.jnilid = jnil.jnilid
		JOIN (
			-- Semester lama: via vw_kelas_tawar
			SELECT vw.id_master_kegiatan, vw.kelas, vw.mkid
			FROM vw_kelas_tawar vw
			WHERE vw.mkid = ? AND vw.kelas = ? AND vw.pakid = ?
			UNION
			-- Semester baru (20252+): via jad
			SELECT jad.id_master_kegiatan, jad.kelas, jad.mkid
			FROM jad
			WHERE jad.mkid = ? AND jad.kelas = ? AND jad.pakid = ?
		) src ON src.id_master_kegiatan = tugas_kul.master_kegiatan_id
		LEFT JOIN mk ON mk.mkid = src.mkid
		LEFT JOIN tugas_submission ts ON ts.tugas_kul_id = tugas_kul.id_tugas_kul AND ts.mhsid = ?
		WHERE (tugas_kul.waktu_mulai_tugas IS NULL OR tugas_kul.waktu_mulai_tugas <= NOW())
		  AND (tugas_kul.waktu_akhir_tugas IS NULL OR tugas_kul.waktu_akhir_tugas >= NOW())
		ORDER BY tugas_kul.id_tugas_kul
	`

	err := s.db.WithContext(ctx).Raw(rawSQL, mkID, class, pakID, mkID, class, pakID, mhsID).Scan(&activeAssignments).Error

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

	rawSQL := `
		SELECT DISTINCT ON (tugas_kul.id_tugas_kul)
			tugas_kul.*,
			mk.mknama,
			src.kelas,
			(ts.id_tugas_submission IS NOT NULL) AS sudah_submit,
			ts.file_tugas AS submission_file,
			ts.link_tugas AS submission_link,
			ts.created_at AS submission_date
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
		LEFT JOIN tugas_submission ts ON ts.tugas_kul_id = tugas_kul.id_tugas_kul AND ts.mhsid = ?
		WHERE src.mhsid = ?
		  AND src.pakid = ?
		  AND (tugas_kul.waktu_mulai_tugas IS NULL OR tugas_kul.waktu_mulai_tugas <= NOW())
		  AND (tugas_kul.waktu_akhir_tugas IS NULL OR tugas_kul.waktu_akhir_tugas >= NOW())
	`

	err := s.db.WithContext(ctx).Raw(rawSQL, mhsID, mhsID, pakID).Scan(&activeAssignments).Error

	return activeAssignments, err
}
