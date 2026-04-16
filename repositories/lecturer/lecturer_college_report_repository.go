package lecturer_repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type lecturerCollegeReportRepository struct {
	db *gorm.DB
}

type LecturerCollegeReportRepository interface {
	GetSubjectCollegeReport(ctx context.Context, mkID string, class string, hourID string, collegeType string) ([]entities.Lecture, error)
	CreateCollege(ctx context.Context, lecture entities.StoreLecture, materials []entities.LectureMaterial) error
	EditCollege(ctx context.Context, lecture entities.Lecture, materials []entities.LectureMaterial) error
	DeleteCollege(ctx context.Context, kulid string) error
	GetSubjectCollegeReportByKulID(ctx context.Context, kulID string) (entities.Lecture, error)
	GetTeamWeeks(ctx context.Context, dosID string, mkID string, kelas string, pakID string) ([]entities.Week, error)
	GetRPSDetail(ctx context.Context, mkID string, weekID string) (map[string]interface{}, error)
}

func NewLecturerCollegeReportRepository(db *gorm.DB) *lecturerCollegeReportRepository {
	return &lecturerCollegeReportRepository{
		db: db,
	}
}

func (r *lecturerCollegeReportRepository) GetSubjectCollegeReport(ctx context.Context, mkID string, class string, hourID string, collegeType string) ([]entities.Lecture, error) {
	lectures := []entities.Lecture{}

	query := `SELECT kul.*, (SELECT COUNT(*) FROM absen WHERE absen.pakid = kul.pakid AND absen.mkid = kul.mkid and absen.jurid = kul.jurid AND absen.kelas = kul.kelas and absen.kultype = kul.kultype and absen.jamid = kul.jamid and absen.weekid = kul.weekid) as mahasiswa_absen FROM kul where kul.mkid = ? and kelas = ? and kultype = ?`
	args := []interface{}{mkID, class, collegeType}

	if hourID != "" {
		query += ` and jamid = ?`
		args = append(args, hourID)
	}

	query += ` order by weekid ASC`

	err := r.db.WithContext(ctx).Raw(query, args...).Find(&lectures).Error

	return lectures, err
}

func (r *lecturerCollegeReportRepository) GetSubjectCollegeReportByKulID(ctx context.Context, kulID string) (entities.Lecture, error) {
	lectures := entities.Lecture{}

	err := r.db.WithContext(ctx).Table("kul").Where("kulid = ?", kulID).Find(&lectures).Error

	return lectures, err
}

func (r *lecturerCollegeReportRepository) CreateCollege(ctx context.Context, lecture entities.StoreLecture, materials []entities.LectureMaterial) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.WithContext(ctx).Create(&lecture).Error

			if err != nil {
				return err
			}

			if len(materials) > 0 {
				err = tx.WithContext(ctx).Create(&materials).Error
				return err
			}

			return nil
		},
	)
}

func (r *lecturerCollegeReportRepository) EditCollege(ctx context.Context, lecture entities.Lecture, materials []entities.LectureMaterial) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.WithContext(ctx).Table("kul").Where("kulid = ?", lecture.LectureID).Updates(&lecture).Error

			if err != nil {
				return err
			}

			err = tx.WithContext(ctx).Where("kul_id = ?", lecture.LectureID).Delete(&entities.LectureMaterial{}).Error

			if err != nil {
				return err
			}

			if len(materials) > 0 {
				err = tx.WithContext(ctx).Model(&entities.LectureMaterial{}).Create(&materials).Error
				return err
			}

			return nil
		},
	)
}

func (r *lecturerCollegeReportRepository) DeleteCollege(ctx context.Context, kulid string) error {
	return r.db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.WithContext(ctx).Where("kul_id = ?", kulid).Delete(&entities.LectureMaterial{}).Error

			if err != nil {
				return err
			}

			err = tx.WithContext(ctx).Exec("DELETE FROM kul where kulid = ?", kulid).Error

			if err != nil {
				return err
			}

			return nil
		},
	)
}

func (r *lecturerCollegeReportRepository) GetTeamWeeks(ctx context.Context, dosID string, mkID string, kelas string, pakID string) ([]entities.Week, error) {
	weeks := []entities.Week{}

	// Cek apakah dosen adalah anggota team teaching
	// Kita join jad untuk mendapatkan id_master_kegiatan karena parameter yang dimiliki adalah mkid dan kelas
	err := r.db.WithContext(ctx).Raw(`
		SELECT week.* FROM jadteamweek
		JOIN week ON week.weekid = jadteamweek.weekid
		JOIN jad ON jad.id_master_kegiatan = jadteamweek.id_jad_master
		WHERE jadteamweek.dosid = ? AND jad.mkid = ? AND jad.kelas = ? AND jad.pakid = ?
		ORDER BY jadteamweek.weekid ASC
	`, dosID, mkID, kelas, pakID).Scan(&weeks).Error

	if err != nil {
		return nil, err
	}

	// Jika tidak ada data spesifik team teaching, kembalikan 16 minggu standar
	if len(weeks) == 0 {
		err = r.db.WithContext(ctx).Table("week").Order("weekid ASC").Find(&weeks).Error
	}

	return weeks, err
}

func (r *lecturerCollegeReportRepository) GetRPSDetail(ctx context.Context, mkID string, weekID string) (map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.WithContext(ctx).Raw(`
		SELECT week.weekid, week.weekno, c.deskripsi_cp, r.deskripsi_rp, r.kode, m.id_mapping_rps 
		FROM week 
		LEFT JOIN mapping_rps as m ON week.weekid = m.weekid 
		JOIN rencana_pembelajaran as r ON m.id_rencana_pembelajaran = r.id_rencana_pembelajaran 
		JOIN capaian_pembelajaran as c ON c.id_capaian_pembelajaran = r.capaian_pembelajaran_id 
		WHERE r.mkid = ? AND week.weekid = ? 
		LIMIT 1;
	`, mkID, weekID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]interface{}{}, nil
	}

	return results[0], nil
}
