package lecturer_services

import (
	"classroom_itats_api/entities"
	lecturer_repositories "classroom_itats_api/repositories/lecturer"
	"context"

	"github.com/google/uuid"
)

type lecturerCollegeReportService struct {
	lecturerCollegeReportRepository lecturer_repositories.LecturerCollegeReportRepository
}

type LecturerCollegeReportService interface {
	GetSubjectCollegeReport(ctx context.Context, mkID string, class string, hourID string, collegeType string) ([]entities.Lecture, error)
	CreateCollege(ctx context.Context, lecture entities.StoreLecture, materials []entities.LectureMaterial) error
	EditCollege(ctx context.Context, lecture entities.Lecture, materials []entities.LectureMaterial) error
	DeleteCollege(ctx context.Context, kulid string) error
	GetSubjectCollegeReportByKulID(ctx context.Context, kulID string) (entities.Lecture, error)
	GetTeamWeeks(ctx context.Context, dosID string, mkID string, kelas string, pakID string) ([]entities.Week, error)
	GetRPSDetail(ctx context.Context, mkID string, weekID string) (map[string]interface{}, error)
}

func NewLecturerCollegeReportService(lecturerCollegeReportRepository lecturer_repositories.LecturerCollegeReportRepository) *lecturerCollegeReportService {
	return &lecturerCollegeReportService{
		lecturerCollegeReportRepository: lecturerCollegeReportRepository,
	}
}

func (r *lecturerCollegeReportService) GetSubjectCollegeReport(ctx context.Context, mkID string, class string, hourID string, collegeType string) ([]entities.Lecture, error) {
	return r.lecturerCollegeReportRepository.GetSubjectCollegeReport(ctx, mkID, class, hourID, collegeType)
}

func (r *lecturerCollegeReportService) CreateCollege(ctx context.Context, lecture entities.StoreLecture, materials []entities.LectureMaterial) error {
	newID, err := uuid.NewRandom()

	if err != nil {
		return err
	}

	lecture.LectureID = newID.String()

	for k, v := range materials {
		newUID, _ := uuid.NewRandom()
		v.LectureMaterialID = newUID.String()
		v.LectID = lecture.LectureID

		materials[k] = v
	}

	return r.lecturerCollegeReportRepository.CreateCollege(ctx, lecture, materials)
}

func (r *lecturerCollegeReportService) EditCollege(ctx context.Context, lecture entities.Lecture, materials []entities.LectureMaterial) error {
	for k, v := range materials {
		newUID, _ := uuid.NewRandom()
		v.LectureMaterialID = newUID.String()
		v.LectID = lecture.LectureID

		materials[k] = v
	}

	return r.lecturerCollegeReportRepository.EditCollege(ctx, lecture, materials)
}

func (r *lecturerCollegeReportService) DeleteCollege(ctx context.Context, kulid string) error {
	return r.lecturerCollegeReportRepository.DeleteCollege(ctx, kulid)
}

func (r *lecturerCollegeReportService) GetSubjectCollegeReportByKulID(ctx context.Context, kulID string) (entities.Lecture, error) {
	return r.lecturerCollegeReportRepository.GetSubjectCollegeReportByKulID(ctx, kulID)
}

func (r *lecturerCollegeReportService) GetTeamWeeks(ctx context.Context, dosID string, mkID string, kelas string, pakID string) ([]entities.Week, error) {
	return r.lecturerCollegeReportRepository.GetTeamWeeks(ctx, dosID, mkID, kelas, pakID)
}

func (r *lecturerCollegeReportService) GetRPSDetail(ctx context.Context, mkID string, weekID string) (map[string]interface{}, error) {
	return r.lecturerCollegeReportRepository.GetRPSDetail(ctx, mkID, weekID)
}
