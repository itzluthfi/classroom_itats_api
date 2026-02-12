package lecturer_services

import (
	"classroom_itats_api/entities"
	lecturer_repositories "classroom_itats_api/repositories/lecturer"
	"context"
)

type lecturerMaterialService struct {
	lecturerMaterialRepository lecturer_repositories.LecturerMaterialRepository
}

type LecturerMaterialService interface {
	GetMaterials(ctx context.Context, dosID string) ([]entities.Material, error)
	GetMaterialSelected(ctx context.Context, kulID string) ([]entities.Material, error)
}

func NewLecturerMaterialService(lecturerMaterialRepository lecturer_repositories.LecturerMaterialRepository) *lecturerMaterialService {
	return &lecturerMaterialService{
		lecturerMaterialRepository: lecturerMaterialRepository,
	}
}

func (l *lecturerMaterialService) GetMaterials(ctx context.Context, dosID string) ([]entities.Material, error) {
	return l.lecturerMaterialRepository.GetMaterials(ctx, dosID)
}

func (l *lecturerMaterialService) GetMaterialSelected(ctx context.Context, kulID string) ([]entities.Material, error) {
	return l.lecturerMaterialRepository.GetMaterialSelected(ctx, kulID)
}
