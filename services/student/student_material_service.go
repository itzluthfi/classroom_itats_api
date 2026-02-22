package student_services

import (
	"classroom_itats_api/entities"
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)

type studentMaterialService struct {
	studentMaterialRepository student_repositories.StudentMaterialRepository
}

type StudentMaterialService interface {
	GetWeekMaterial(ctx context.Context, pakID string, mkID string, class string) ([]entities.LectureWeek, error)
	GetStudyAchievement(ctx context.Context, pakID string, mkID string, class string) ([]entities.StudyAchievement, error)
	GetStudentMaterial(ctx context.Context, mkID string, class string, pakID string, weekID int) ([]entities.Material, error)
	GetStudentAssignment(ctx context.Context, masterActivityID string, weekID float64, mhsID string) ([]entities.Assignment, error)
	GetStudentAssignmentGroup(ctx context.Context, masterActivityID string) ([]entities.Assignment, error)
	GetStudentScore(ctx context.Context, mhsID string, masterActivityID string) ([]entities.StudentAssignmentScore, error)
	GetStudentAssignmentSubmission(ctx context.Context, mhsID string, assignmentID int) (entities.AssignmentSubmission, error)
	GetActiveAssignment(ctx context.Context, masterActivityID string, mhsID string) ([]entities.Assignment, error)
	GetHomeActiveAssignment(ctx context.Context, mhsID string) ([]entities.Assignment, error)
}

func NewStudentMaterialService(studentMaterialRepository student_repositories.StudentMaterialRepository) *studentMaterialService {
	return &studentMaterialService{studentMaterialRepository: studentMaterialRepository}
}

func (s *studentMaterialService) GetWeekMaterial(ctx context.Context, pakID string, mkID string, class string) ([]entities.LectureWeek, error) {
	return s.studentMaterialRepository.GetWeekMaterial(ctx, pakID, mkID, class)
}

func (s *studentMaterialService) GetStudyAchievement(ctx context.Context, pakID string, mkID string, class string) ([]entities.StudyAchievement, error) {
	return s.studentMaterialRepository.GetStudyAchievement(ctx, mkID, class, pakID)
}

func (s *studentMaterialService) GetStudentMaterial(ctx context.Context, mkID string, class string, pakID string, weekID int) ([]entities.Material, error) {
	return s.studentMaterialRepository.GetStudentMaterial(ctx, mkID, class, pakID, weekID)
}

func (s *studentMaterialService) GetStudentAssignment(ctx context.Context, masterActivityID string, weekID float64, mhsID string) ([]entities.Assignment, error) {
	return s.studentMaterialRepository.GetStudentAssignment(ctx, masterActivityID, weekID, mhsID)
}

func (s *studentMaterialService) GetStudentAssignmentGroup(ctx context.Context, masterActivityID string) ([]entities.Assignment, error) {
	return s.studentMaterialRepository.GetStudentAssignmentGroup(ctx, masterActivityID)
}

func (s *studentMaterialService) GetStudentAssignmentSubmission(ctx context.Context, mhsID string, assignmentID int) (entities.AssignmentSubmission, error) {
	return s.studentMaterialRepository.GetStudentAssignmentSubmission(ctx, mhsID, assignmentID)
}

func (s *studentMaterialService) GetStudentScore(ctx context.Context, mhsID string, masterActivityID string) ([]entities.StudentAssignmentScore, error) {
	return s.studentMaterialRepository.GetStudentScore(ctx, mhsID, masterActivityID)
}

func (s *studentMaterialService) GetActiveAssignment(ctx context.Context, masterActivityID string, mhsID string) ([]entities.Assignment, error) {
	return s.studentMaterialRepository.GetActiveAssignment(ctx, masterActivityID, mhsID)
}

func (s *studentMaterialService) GetHomeActiveAssignment(ctx context.Context, mhsID string) ([]entities.Assignment, error) {
	return s.studentMaterialRepository.GetHomeActiveAssignment(ctx, mhsID)
}
