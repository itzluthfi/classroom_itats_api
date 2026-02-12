package lecturer_repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type lecturerMaterialRepository struct {
	db *gorm.DB
}

type LecturerMaterialRepository interface {
	GetMaterials(ctx context.Context, dosID string) ([]entities.Material, error)
	GetMaterialSelected(ctx context.Context, kulID string) ([]entities.Material, error)
}

func NewLecturerMaterialRepository(db *gorm.DB) *lecturerMaterialRepository {
	return &lecturerMaterialRepository{
		db: db,
	}
}

func (l *lecturerMaterialRepository) GetMaterials(ctx context.Context, dosID string) ([]entities.Material, error) {
	materials := []entities.Material{}

	err := l.db.WithContext(ctx).Table("materi").Where("dos_id = ?", dosID).Where("hidden_status = ?", 0).Where("deleted_at is null").Find(&materials).Error

	return materials, err
}

func (l *lecturerMaterialRepository) GetMaterialSelected(ctx context.Context, kulID string) ([]entities.Material, error) {
	materials := []entities.Material{}

	err := l.db.WithContext(ctx).Table("materi").Select(
		"materi.materi_id",
		"kul.dosid",
		"materi.judul_materi",
		"materi.link_materi",
		"materi.created_at",
		"materi.updated_at",
		"materi.deleted_at",
		"materi.hidden_status",
		"kul_materi.kul_materi_id",
		"kul_materi.kul_id",
	).Joins("JOIN kul_materi on kul_materi.materi_id = materi.materi_id").Joins("JOIN kul on kul_materi.kul_id = kul.kulid").Where("kul_id", kulID).Where("hidden_status = ?", 0).Where("deleted_at is null").Find(&materials).Error

	return materials, err
}
