package student_services

import (
	"classroom_itats_api/entities"
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)

type studentPresenceService struct {
	studentPresenceRepository student_repositories.StudentPresenceRepository
}

type StudentPresenceService interface {
	GetStudentPresences(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Presence, error)
	GetPresenceQuestion(ctx context.Context, pakID string) ([]entities.PresenceQuestion, error)
	GetSubjectResponsi(ctx context.Context, pakID string, mkID string, class string) (int, error)
	GetStudentPresencesSeparated(ctx context.Context, pakID string, mkID string, class string, mhsID string) (entities.StudentPresenceJSON, error)
	SetStudentPresenceAnswers(ctx context.Context, PresenceAnswers []entities.PresenceAnswer) error
	SetStudentPresence(ctx context.Context, StudentPresence entities.Presence) error
}

func NewStudentPresenceService(studentPresenceRepository student_repositories.StudentPresenceRepository) *studentPresenceService {
	return &studentPresenceService{studentPresenceRepository: studentPresenceRepository}
}

func (s *studentPresenceService) GetStudentPresences(ctx context.Context, pakID string, mkID string, class string, mhsID string) ([]entities.Presence, error) {
	return s.studentPresenceRepository.GetStudentPresences(ctx, pakID, mkID, class, mhsID)
}

func (s *studentPresenceService) GetPresenceQuestion(ctx context.Context, pakID string) ([]entities.PresenceQuestion, error) {
	return s.studentPresenceRepository.GetPresenceQuestion(ctx, pakID)
}

func (s *studentPresenceService) GetSubjectResponsi(ctx context.Context, pakID string, mkID string, class string) (int, error) {
	return s.studentPresenceRepository.GetSubjectResponsi(ctx, pakID, mkID, class)
}

func (s *studentPresenceService) GetStudentPresencesSeparated(ctx context.Context, pakID string, mkID string, class string, mhsID string) (entities.StudentPresenceJSON, error) {
	presences, err := s.studentPresenceRepository.GetStudentPresences(ctx, pakID, mkID, class, mhsID)

	if err != nil {
		return entities.StudentPresenceJSON{}, err
	}

	materialPresences := []entities.Presence{}
	allResponsiPresences := [][]entities.Presence{}
	responsiPresences := []entities.Presence{}

	for key, v := range presences {
		if key > 0 {
			// if v.CollegeType != "M" {}
			if v.CollegeType != presences[key-1].CollegeType {
				materialPresences = append(materialPresences, presences[:key]...)
				responsiPresences = append(responsiPresences, presences[key:]...)
				break
			}
		}
	}

	temp := 0
	for key := range responsiPresences {
		if key+1 == len(responsiPresences) {
			allResponsiPresences = append(allResponsiPresences, responsiPresences[temp:])
		} else {
			if responsiPresences[key+1].WeekID == 1 {
				allResponsiPresences = append(allResponsiPresences, responsiPresences[:key])
				temp = key
			}
		}
	}

	if len(materialPresences) == 0 {
		return entities.StudentPresenceJSON{
			MaterialPresences: presences,
			ResponsiPresences: allResponsiPresences,
		}, err
	}

	return entities.StudentPresenceJSON{
		MaterialPresences: materialPresences,
		ResponsiPresences: allResponsiPresences,
	}, err
}

func (s *studentPresenceService) SetStudentPresenceAnswers(ctx context.Context, PresenceAnswers []entities.PresenceAnswer) error {
	return s.studentPresenceRepository.SetStudentPresenceAnswers(ctx, PresenceAnswers)
}

func (s *studentPresenceService) SetStudentPresence(ctx context.Context, StudentPresence entities.Presence) error {
	return s.studentPresenceRepository.SetStudentPresence(ctx, StudentPresence)
}
