package student_services

import (
	"classroom_itats_api/entities"
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)

type studentProfileService struct {
	studentProfileRepository student_repositories.StudentProfileRepository
}

type StudentProfileService interface {
	GetStudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfile, error)
	GetStudentSubjectPresence(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubjectPresence, error)
	StudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfileJSON, error)
}

func NewStudentProfileService(studentProfileRepository student_repositories.StudentProfileRepository) *studentProfileService {
	return &studentProfileService{studentProfileRepository: studentProfileRepository}
}

func (s *studentProfileService) GetStudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfile, error) {
	return s.studentProfileRepository.GetStudentProfile(ctx, mhsID, pakID)
}

func (s *studentProfileService) GetStudentSubjectPresence(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubjectPresence, error) {
	return s.studentProfileRepository.GetStudentSubjectPresence(ctx, mhsID, pakID)
}

func (s *studentProfileService) StudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfileJSON, error) {
	studentProfileJSON := entities.StudentProfileJSON{}

	studentProfile, err := s.studentProfileRepository.GetStudentProfile(ctx, mhsID, pakID)

	if err != nil {
		return entities.StudentProfileJSON{}, err
	}

	studentSubjectPresence, err := s.studentProfileRepository.GetStudentSubjectPresence(ctx, mhsID, pakID)

	if err != nil {
		return entities.StudentProfileJSON{}, err
	}

	studentProfileJSON = entities.StudentProfileJSON{
		UserID:                  studentProfile.UserID,
		Name:                    studentProfile.Name,
		PhoneNumber:             studentProfile.PhoneNumber,
		Email:                   studentProfile.Email,
		Photo:                   studentProfile.Photo,
		Presence:                studentProfile.Presence,
		TotalPresence:           studentProfile.TotalPresence,
		AssignmentSubmited:      studentProfile.AssignmentSubmited,
		TotalAssignment:         studentProfile.TotalAssignment,
		StudentSubjectPresences: studentSubjectPresence,
	}

	return studentProfileJSON, err
}
