package student_services

import (
	"classroom_itats_api/entities"
	student_repositories "classroom_itats_api/repositories/student"
	"context"
	"fmt"
	"strings"
	"time"
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

// computeClassStatus computes the class status for each schedule based on
// active room key loans and current time/day
func computeClassStatus(schedules []entities.SubjectSchedule, loans []entities.RoomKeyLoan, todayHariCode string, now time.Time) []entities.SubjectSchedule {
	// Build a set of room IDs that currently have active loans
	activeLoanRooms := make(map[string]bool)
	for _, loan := range loans {
		roomID := strings.TrimSpace(loan.RoomID)
		if roomID != "" {
			activeLoanRooms[roomID] = true
		}
	}

	currentTime := now.Format("15:04:05")

	for i := range schedules {
		schedule := &schedules[i]
		dayCode := strings.TrimSpace(schedule.DayCode)
		roomID := strings.TrimSpace(schedule.SubjectRoom)

		isToday := todayHariCode != "" && dayCode == todayHariCode
		hasLoan := roomID != "" && activeLoanRooms[roomID]

		if !isToday {
			schedule.ClassStatus = ""
			continue
		}

		// Calculate end time based on jammulai + SKS * 50 minutes
		timeStart := strings.TrimSpace(schedule.TimeStart)
		inSchedule := false

		if timeStart != "" {
			// Parse jammulai as HH:MM:SS or HH:MM
			startParsed, err := time.Parse("15:04:05", timeStart)
			if err != nil {
				startParsed, err = time.Parse("15:04", timeStart)
			}
			if err == nil {
				sks := schedule.SKS
				if sks <= 0 {
					sks = 2 // default fallback
				}
				endParsed := startParsed.Add(time.Duration(sks*50) * time.Minute)
				endTimeStr := endParsed.Format("15:04:05")

				inSchedule = currentTime >= startParsed.Format("15:04:05") && currentTime <= endTimeStr
			}
		}

		if hasLoan && inSchedule {
			schedule.ClassStatus = "sedang_berlangsung"
		} else if hasLoan {
			schedule.ClassStatus = "kunci_diambil"
		} else if inSchedule {
			schedule.ClassStatus = "jadwal_aktif"
		} else {
			schedule.ClassStatus = ""
		}
	}

	return schedules
}

// fetchClassStatusData fetches active loans and today's hari code for status computation
func (s *studentSubjectService) fetchClassStatusData(ctx context.Context) ([]entities.RoomKeyLoan, string, time.Time) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)

	// dayOfWeekIso: Monday=1 .. Sunday=7
	dayOfWeek := int(now.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7 // Sunday = 7
	}

	loans, err := s.studentSubjectRepository.GetActiveRoomLoans(ctx)
	if err != nil {
		fmt.Println("Warning: failed to fetch active room loans:", err)
		loans = []entities.RoomKeyLoan{}
	}

	todayHariCode, err := s.studentSubjectRepository.GetTodayHariCode(ctx, dayOfWeek)
	if err != nil {
		fmt.Println("Warning: failed to fetch today hari code:", err)
		todayHariCode = ""
	}

	return loans, todayHariCode, now
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

	// Fetch class status data and compute status
	loans, todayHariCode, now := s.fetchClassStatusData(ctx)
	subjectSchedules = computeClassStatus(subjectSchedules, loans, todayHariCode, now)

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

	// Fetch class status data and compute status
	loans, todayHariCode, now := s.fetchClassStatusData(ctx)
	subjectSchedules = computeClassStatus(subjectSchedules, loans, todayHariCode, now)

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

