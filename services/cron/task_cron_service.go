package cron_service

import (
	student_repositories "classroom_itats_api/repositories/student"
	"context"
)
//tes
type taskCronService struct {
	studentMaterialRepository student_repositories.StudentMaterialRepository
}

type TaskCronService interface {
	AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error)
	AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error)
	AssignmentReminderH1(ctx context.Context) ([]map[string]interface{}, error)
}

func NewTaskCronService(studentMaterialRepository student_repositories.StudentMaterialRepository) *taskCronService {
	return &taskCronService{
		studentMaterialRepository: studentMaterialRepository,
	}
}

func (t *taskCronService) AssignmentCreated(ctx context.Context) ([]map[string]interface{}, error) {
	return t.studentMaterialRepository.AssignmentCreated(ctx)
}

func (t *taskCronService) AssignmentReminder(ctx context.Context) ([]map[string]interface{}, error) {
	return t.studentMaterialRepository.AssignmentReminder(ctx)
}

func (t *taskCronService) AssignmentReminderH1(ctx context.Context) ([]map[string]interface{}, error) {
	return t.studentMaterialRepository.AssignmentReminderH1(ctx)
}
