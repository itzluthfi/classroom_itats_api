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
