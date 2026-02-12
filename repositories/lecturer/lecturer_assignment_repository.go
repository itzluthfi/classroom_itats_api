package lecturer_repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type lecturerAssignmentRepository struct {
	db *gorm.DB
}

type LecturerAssignmentRepository interface {
	GetLecturerCreatedAssignment(ctx context.Context, pakID string, dosID string) ([]entities.Assignment, error)
	GetWeekAssignment(ctx context.Context) ([]entities.Week, error)
	GetScoreTypeAssignment(ctx context.Context) ([]entities.ScoreType, error)
}

func NewLecturerAssignmentRepository(db *gorm.DB) *lecturerAssignmentRepository {
	return &lecturerAssignmentRepository{
		db: db,
	}
}

func (l *lecturerAssignmentRepository) GetLecturerCreatedAssignment(ctx context.Context, pakID string, dosID string) ([]entities.Assignment, error) {
	assignments := []entities.Assignment{}

	err := l.db.WithContext(ctx).Raw(`
		select tugas_kul.*, klskul.kelas, klskul.mknama,jnil.jnildesc, (SELECT count(*) FROM tugas_submission where tugas_kul_id = tugas_kul.id_tugas_kul) AS jml_pengumpulan
		from tugas_kul
		join (SELECT DISTINCT id_master_kegiatan, mknama, kelas, dosid, pakid FROM vw_kelas_tawar) AS klskul on
		tugas_kul.master_kegiatan_id = klskul.id_master_kegiatan
		join jnil on jnil.jnilid = tugas_kul.jnilid
		where (tugas_kul.master_kegiatan_id in ((select id_master_kegiatan from vw_kelas_tawar where pakid = ? and dosid = ?))
		or tugas_kul.master_kegiatan_id in ((select id_jad_master from jadteam where dosid = ?)))
	`, pakID, dosID, dosID).Find(&assignments).Error

	return assignments, err
}

func (l *lecturerAssignmentRepository) GetWeekAssignment(ctx context.Context) ([]entities.Week, error) {
	weeks := []entities.Week{}

	err := l.db.WithContext(ctx).Table("week").Order("weekid ASC").Find(&weeks).Error

	return weeks, err
}

func (l *lecturerAssignmentRepository) GetScoreTypeAssignment(ctx context.Context) ([]entities.ScoreType, error) {
	scoreTypes := []entities.ScoreType{}

	err := l.db.WithContext(ctx).Table("jnil").Where("jnilaktif = ?", true).Find(&scoreTypes).Error

	return scoreTypes, err
}
