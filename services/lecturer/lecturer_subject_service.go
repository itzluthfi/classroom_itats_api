package lecturer_services

import (
	"classroom_itats_api/entities"
	lecturer_repositories "classroom_itats_api/repositories/lecturer"
	"context"
	"fmt"
)

type lecturerSubjectService struct {
	lecturerSubjectRepository lecturer_repositories.LecturerSubjectRepository
}

type LecturerSubjectService interface {
	GetActiveSubjectByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubject, error)
	LecturePeriodes(ctx context.Context, dosID string) ([]entities.AcademicPeriod, error)
	GetLecturerSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubject, error)
	LecturerSubjectMajor(ctx context.Context, dosID string, pakID string) ([]entities.Major, error)
	GetSubjectClassLecturer(ctx context.Context, dosID string) ([]entities.SubjectSchedule, error)
	GetSubjectClassLecturerFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectSchedule, error)
	GetActiveLecturerSubject(ctx context.Context, dosID string) ([]entities.SubjectJSON, error)
	GetFilteredLecturerSubject(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectJSON, error)
	GetStudentScore(ctx context.Context, pakid string, mkid string, class string) ([]entities.StudentScore, error)
	GetAssignment(ctx context.Context, masterActivityID string) ([]entities.Assignment, error)
	GetAssignmentSubmission(ctx context.Context, masterActivityID string) ([]entities.AssignmentSubmission, error)
	GetClassOffered(ctx context.Context, masterActivityID string) (entities.ClassOffered, error)
	GetSubjectPercentageScore(ctx context.Context, masterActivityID string) (entities.PercentageScore, error)
	GetActiveSubjectReportByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubjectReport, error)
	GetLecturerSubjectReportFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubjectReport, error)
}

func NewLecturerSubjectService(lecturerSubjectRepository lecturer_repositories.LecturerSubjectRepository) *lecturerSubjectService {
	return &lecturerSubjectService{lecturerSubjectRepository: lecturerSubjectRepository}
}

func (l *lecturerSubjectService) GetActiveSubjectByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubject, error) {
	return l.lecturerSubjectRepository.GetActiveSubjectByLecturerID(ctx, dosID)
}

func (l *lecturerSubjectService) GetSubjectClassLecturer(ctx context.Context, dosID string) ([]entities.SubjectSchedule, error) {
	return l.lecturerSubjectRepository.GetSubjectClassLecturer(ctx, dosID)
}

func (l *lecturerSubjectService) GetActiveLecturerSubject(ctx context.Context, dosID string) ([]entities.SubjectJSON, error) {
	activeLecturerSubjects := []entities.SubjectJSON{}

	lecturerSubjects, err := l.GetActiveSubjectByLecturerID(ctx, dosID)
	if err != nil {
		return nil, err
	}

	subjectSchedules, err := l.GetSubjectClassLecturer(ctx, dosID)
	if err != nil {
		return nil, err
	}

	for _, lecturerSubject := range lecturerSubjects {
		subjectClass := []entities.SubjectSchedule{}
		for _, subjectSchedule := range subjectSchedules {
			if subjectSchedule.SubjectClass == lecturerSubject.SubjectClass && subjectSchedule.SubjectID == lecturerSubject.SubjectID {
				subjectClass = append(subjectClass, subjectSchedule)
			}
		}
		activeLecturerSubjects = append(activeLecturerSubjects, entities.SubjectJSON{
			SubjectClass:     lecturerSubject.SubjectClass,
			SubjectID:        lecturerSubject.SubjectID,
			SubjectCredits:   lecturerSubject.SubjectCredits,
			MajorID:          lecturerSubject.MajorID,
			MajorName:        lecturerSubject.MajorName,
			AcademicPeriodID: lecturerSubject.AcademicPeriodID,
			LecturerID:       lecturerSubject.LecturerID,
			LecturerName:     lecturerSubject.LecturerName,
			SubjectName:      lecturerSubject.SubjectName,
			TotalStudent:     lecturerSubject.TotalStudent,
			ActivityMasterID: lecturerSubject.ActivityMasterID,
			SubjectSchedules: subjectClass,
		})
	}

	return activeLecturerSubjects, nil
}

func (l *lecturerSubjectService) GetLecturerSubjectFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubject, error) {
	if jurID == "" {
		return l.lecturerSubjectRepository.GetSubjectByLecturerIDWithPeriod(ctx, dosID, pakID)
	}
	return l.lecturerSubjectRepository.GetLecturerSubjectFiltered(ctx, dosID, pakID, jurID)
}

func (l *lecturerSubjectService) GetSubjectClassLecturerFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectSchedule, error) {
	if jurID == "" {
		return l.lecturerSubjectRepository.GetSubjectClassLecturerWithPeriod(ctx, dosID, pakID)
	}
	return l.lecturerSubjectRepository.GetSubjectClassLecturerFiltered(ctx, dosID, pakID, jurID)
}

func (l *lecturerSubjectService) GetFilteredLecturerSubject(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.SubjectJSON, error) {
	activeLecturerSubjects := []entities.SubjectJSON{}

	lecturerSubjects, err := l.GetLecturerSubjectFiltered(ctx, dosID, pakID, jurID)
	if err != nil {
		return nil, err
	}

	subjectSchedules, err := l.GetSubjectClassLecturerFiltered(ctx, dosID, pakID, jurID)
	if err != nil {
		return nil, err
	}

	for _, lecturerSubject := range lecturerSubjects {
		subjectClass := []entities.SubjectSchedule{}
		for _, subjectSchedule := range subjectSchedules {
			if subjectSchedule.SubjectClass == lecturerSubject.SubjectClass && subjectSchedule.SubjectID == lecturerSubject.SubjectID {
				subjectClass = append(subjectClass, subjectSchedule)
			}
		}
		activeLecturerSubjects = append(activeLecturerSubjects, entities.SubjectJSON{
			SubjectClass:     lecturerSubject.SubjectClass,
			SubjectID:        lecturerSubject.SubjectID,
			SubjectCredits:   lecturerSubject.SubjectCredits,
			MajorID:          lecturerSubject.MajorID,
			MajorName:        lecturerSubject.MajorName,
			AcademicPeriodID: lecturerSubject.AcademicPeriodID,
			LecturerID:       lecturerSubject.LecturerID,
			LecturerName:     lecturerSubject.LecturerName,
			SubjectName:      lecturerSubject.SubjectName,
			TotalStudent:     lecturerSubject.TotalStudent,
			ActivityMasterID: lecturerSubject.ActivityMasterID,
			SubjectSchedules: subjectClass,
		})
	}

	return activeLecturerSubjects, nil
}

func (l *lecturerSubjectService) LecturePeriodes(ctx context.Context, dosID string) ([]entities.AcademicPeriod, error) {
	return l.lecturerSubjectRepository.LecturePeriodes(ctx, dosID)
}

func (l *lecturerSubjectService) LecturerSubjectMajor(ctx context.Context, dosID string, pakID string) ([]entities.Major, error) {
	if pakID == "" {
		academicPeriods, err := l.LecturePeriodes(ctx, dosID)

		if err != nil {
			return []entities.Major{}, err
		}

		for _, v := range academicPeriods {
			if v.IsActive {
				pakID = v.AcademicPeriodID
			}
		}
	}
	return l.lecturerSubjectRepository.LecturerSubjectMajor(ctx, dosID, pakID)
}

func (l *lecturerSubjectService) GetStudentScore(ctx context.Context, pakid string, mkid string, class string) ([]entities.StudentScore, error) {
	return l.lecturerSubjectRepository.GetStudentScore(ctx, pakid, mkid, class)
}

func (l *lecturerSubjectService) GetAssignment(ctx context.Context, masterActivityID string) ([]entities.Assignment, error) {
	return l.lecturerSubjectRepository.GetAssignment(ctx, masterActivityID)
}

func (l *lecturerSubjectService) GetAssignmentSubmission(ctx context.Context, masterActivityID string) ([]entities.AssignmentSubmission, error) {
	return l.lecturerSubjectRepository.GetAssignmentSubmission(ctx, masterActivityID)
}

func (l *lecturerSubjectService) GetClassOffered(ctx context.Context, masterActivityID string) (entities.ClassOffered, error) {
	return l.lecturerSubjectRepository.GetClassOffered(ctx, masterActivityID)
}

func (l *lecturerSubjectService) GetSubjectPercentageScore(ctx context.Context, masterActivityID string) (entities.PercentageScore, error) {
	assignments, err := l.GetAssignment(ctx, masterActivityID)
	if err != nil {
		return entities.PercentageScore{}, err
	}

	// assignmentSubmissons, err := l.GetAssignmentSubmission(ctx, masterActivityID)
	// if err != nil {
	// 	return nil, err
	// }

	classOffered, err := l.GetClassOffered(ctx, masterActivityID)
	if err != nil {
		return entities.PercentageScore{}, err
	}

	detail := []entities.PercentageScoreDetail{}

	detail = append(detail, entities.PercentageScoreDetail{
		ID:              classOffered.SubjectID,
		AssignmentTitle: "Presensi",
		AssignmentType:  "-",
		WeekID:          "-",
		Percentage:      classOffered.PresencePercentageScore,
	})

	for _, v := range assignments {
		detail = append(detail, entities.PercentageScoreDetail{
			ID:              fmt.Sprintf("%d", v.AssignmentID),
			AssignmentTitle: v.AssignmentTitle,
			AssignmentType:  v.JnilDesc,
			WeekID:          fmt.Sprintf("%d", v.WeekID),
			Percentage:      v.RealPrercentage,
		})
	}

	totalPercentage := 0.0
	for _, v := range detail {
		totalPercentage += v.Percentage
	}

	data := entities.PercentageScore{
		TotalPercentage:        totalPercentage,
		PercentageScoreDetails: detail,
	}

	return data, err
}

func (l *lecturerSubjectService) GetActiveSubjectReportByLecturerID(ctx context.Context, dosID string) ([]entities.LecturerSubjectReport, error) {
	return l.lecturerSubjectRepository.GetActiveSubjectReportByLecturerID(ctx, dosID)
}

func (l *lecturerSubjectService) GetLecturerSubjectReportFiltered(ctx context.Context, dosID string, pakID string, jurID string) ([]entities.LecturerSubjectReport, error) {
	if jurID == "" {
		return l.lecturerSubjectRepository.GetSubjectReportByLecturerIDWithPeriod(ctx, dosID, pakID)
	}
	return l.lecturerSubjectRepository.GetLecturerReportSubjectFiltered(ctx, dosID, pakID, jurID)
}
