package student_services

import (
	"classroom_itats_api/entities"
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)

type studentSubjectService struct {
	studentSubjectRepository student_repositories.StudentSubjectRepository
}

type StudentSubjectService interface {
	GetActiveSubjectByStudentID(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error)
	GetSubjectClassStudent(ctx context.Context, mhsID string) ([]entities.SubjectSchedule, error)
	GetActiveStudentSubject(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectJSON, error)
	GetSubjectByStudentIDWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error)
	GetSubjectClassStudentWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectSchedule, error)
	GetStudentSubjectByPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectJSON, error)
	StudyPeriodes(ctx context.Context, mhsID string) ([]entities.AcademicPeriod, error)
}

func NewStudentSubjectService(studentSubjectRepository student_repositories.StudentSubjectRepository) *studentSubjectService {
	return &studentSubjectService{
		studentSubjectRepository: studentSubjectRepository,
	}
}

func (s *studentSubjectService) GetActiveSubjectByStudentID(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error) {
	return s.studentSubjectRepository.GetActiveSubjectByStudentID(ctx, mhsID, pakID)
}

func (s *studentSubjectService) GetSubjectClassStudent(ctx context.Context, mhsID string) ([]entities.SubjectSchedule, error) {
	return s.studentSubjectRepository.GetSubjectClassStudent(ctx, mhsID)
}

func (s *studentSubjectService) GetActiveStudentSubject(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectJSON, error) {
	activeStudentSubject := []entities.SubjectJSON{}

	studentSubjects, err := s.GetActiveSubjectByStudentID(ctx, mhsID, pakID)
	if err != nil {
		return nil, err
	}

	var subjectSchedules []entities.SubjectSchedule
	if pakID != "" {
		subjectSchedules, err = s.GetSubjectClassStudentWithPeriod(ctx, mhsID, pakID)
	} else {
		subjectSchedules, err = s.GetSubjectClassStudent(ctx, mhsID)
	}

	if err != nil {
		return nil, err
	}

	for _, studentSubject := range studentSubjects {
		subjectClass := []entities.SubjectSchedule{}
		for _, subjectSchedule := range subjectSchedules {
			if subjectSchedule.SubjectID == studentSubject.SubjectID {
				subjectClass = append(subjectClass, subjectSchedule)
			}
		}
		activeStudentSubject = append(activeStudentSubject, entities.SubjectJSON{
			SubjectClass:     studentSubject.SubjectClass,
			SubjectID:        studentSubject.SubjectID,
			SubjectCredits:   studentSubject.SubjectCredits,
			MajorID:          studentSubject.MajorID,
			AcademicPeriodID: studentSubject.AcademicPeriodID,
			LecturerName:     studentSubject.LecturerName,
			LecturerID:       studentSubject.LecturerID,
			SubjectName:      studentSubject.SubjectName,
			SubjectSchedules: subjectClass,
			ActivityMasterID: studentSubject.ActivityMasterID,
		})
	}

	return activeStudentSubject, nil
}

func (s *studentSubjectService) GetSubjectByStudentIDWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubject, error) {
	return s.studentSubjectRepository.GetSubjectByStudentIDWithPeriod(ctx, mhsID, pakID)
}

func (s *studentSubjectService) GetSubjectClassStudentWithPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectSchedule, error) {
	return s.studentSubjectRepository.GetSubjectClassStudentWithPeriod(ctx, mhsID, pakID)
}

func (s *studentSubjectService) GetStudentSubjectByPeriod(ctx context.Context, mhsID string, pakID string) ([]entities.SubjectJSON, error) {
	activeStudentSubjects := []entities.SubjectJSON{}

	studentSubjects, err := s.GetSubjectByStudentIDWithPeriod(ctx, mhsID, pakID)
	if err != nil {
		return nil, err
	}

	subjectSchedules, err := s.GetSubjectClassStudentWithPeriod(ctx, mhsID, pakID)
	if err != nil {
		return nil, err
	}

	for _, studentSubject := range studentSubjects {
		subjectClass := []entities.SubjectSchedule{}
		for _, subjectSchedule := range subjectSchedules {
			if subjectSchedule.SubjectID == studentSubject.SubjectID {
				subjectClass = append(subjectClass, subjectSchedule)
			}
		}
		activeStudentSubjects = append(activeStudentSubjects, entities.SubjectJSON{
			SubjectClass:     studentSubject.SubjectClass,
			SubjectID:        studentSubject.SubjectID,
			SubjectCredits:   studentSubject.SubjectCredits,
			MajorID:          studentSubject.MajorID,
			AcademicPeriodID: studentSubject.AcademicPeriodID,
			LecturerName:     studentSubject.LecturerName,
			LecturerID:       studentSubject.LecturerID,
			SubjectName:      studentSubject.SubjectName,
			SubjectSchedules: subjectClass,
			ActivityMasterID: studentSubject.ActivityMasterID,
		})
	}

	return activeStudentSubjects, nil
}

func (s *studentSubjectService) StudyPeriodes(ctx context.Context, mhsID string) ([]entities.AcademicPeriod, error) {
	return s.studentSubjectRepository.StudyPeriodes(ctx, mhsID)
}
