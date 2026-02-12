package cron_service

import (
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)

type presenceCronService struct {
	studentPresenceRepository student_repositories.StudentPresenceRepository
}

type PresenceCronService interface {
	PresenceCreated(ctx context.Context) ([]map[string]interface{}, error)
	PresenceReminder(ctx context.Context) ([]map[string]interface{}, error)
}

func NewPresenceCronService(studentPresenceRepository student_repositories.StudentPresenceRepository) *presenceCronService {
	return &presenceCronService{studentPresenceRepository: studentPresenceRepository}
}

func (p *presenceCronService) PresenceCreated(ctx context.Context) ([]map[string]interface{}, error) {
	return p.studentPresenceRepository.PresenceCreated(ctx)
}
func (p *presenceCronService) PresenceReminder(ctx context.Context) ([]map[string]interface{}, error) {
	return p.studentPresenceRepository.PresenceReminder(ctx)
}
