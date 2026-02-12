package lecturer_services

import (
	"classroom_itats_api/entities"
	lecturer_repositories "classroom_itats_api/repositories/lecturer"
	"context"
)

type lecturerAssignmentService struct {
	lecturerAssignmentRepository lecturer_repositories.LecturerAssignmentRepository
}

type LecturerAssignmentService interface {
	GetLecturerCreatedAssignment(ctx context.Context, pakID string, dosID string) ([]entities.Assignment, error)
	GetWeekAssignment(ctx context.Context) ([]entities.Week, error)
	GetScoreTypeAssignment(ctx context.Context) ([]entities.ScoreType, error)
}

func NewLecturerAssignmentService(lecturerAssignmentRepository lecturer_repositories.LecturerAssignmentRepository) *lecturerAssignmentService {
	return &lecturerAssignmentService{
		lecturerAssignmentRepository: lecturerAssignmentRepository,
	}
}

func (l *lecturerAssignmentService) GetLecturerCreatedAssignment(ctx context.Context, pakID string, dosID string) ([]entities.Assignment, error) {
	return l.lecturerAssignmentRepository.GetLecturerCreatedAssignment(ctx, pakID, dosID)
}

func (l *lecturerAssignmentService) GetWeekAssignment(ctx context.Context) ([]entities.Week, error) {
	return l.lecturerAssignmentRepository.GetWeekAssignment(ctx)
}

func (l *lecturerAssignmentService) GetScoreTypeAssignment(ctx context.Context) ([]entities.ScoreType, error) {
	return l.lecturerAssignmentRepository.GetScoreTypeAssignment(ctx)
}
