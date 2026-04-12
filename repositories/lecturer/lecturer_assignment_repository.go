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
		SELECT DISTINCT ON (tugas_kul.id_tugas_kul)
			tugas_kul.*,
			src.kelas,
			mk.mknama,
			jnil.jnildesc,
			(SELECT count(*) FROM tugas_submission WHERE tugas_kul_id = tugas_kul.id_tugas_kul) AS jml_pengumpulan
		FROM tugas_kul
		JOIN (
			-- Jalur lama: via vw_kelas_tawar
			SELECT id_master_kegiatan, kelas, mknama, dosid, pakid
			FROM vw_kelas_tawar
			WHERE pakid = ? AND dosid = ?
			UNION
			-- Jalur baru: via jad
			SELECT jad.id_master_kegiatan, jad.kelas, mk2.mknama, jad.dosid, jad.pakid
			FROM jad
			JOIN mk mk2 ON mk2.mkid = jad.mkid
			WHERE jad.pakid = ? AND jad.dosid = ?
			UNION
			-- Jalur tim dosen: via jadteam
			SELECT jad2.id_master_kegiatan, jad2.kelas, mk3.mknama, jt.dosid, jad2.pakid
			FROM jadteam jt
			JOIN jad jad2 ON jad2.id_master_kegiatan = jt.id_jad_master
			JOIN mk mk3 ON mk3.mkid = jad2.mkid
			WHERE jt.dosid = ?
		) src ON src.id_master_kegiatan = tugas_kul.master_kegiatan_id
		JOIN jnil ON jnil.jnilid = tugas_kul.jnilid
		LEFT JOIN mk ON mk.mkid = COALESCE(
			(SELECT mkid FROM vw_kelas_tawar WHERE id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1),
			(SELECT mkid FROM jad WHERE id_master_kegiatan = tugas_kul.master_kegiatan_id LIMIT 1)
		)
	`, pakID, dosID, pakID, dosID, dosID).Scan(&assignments).Error

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
