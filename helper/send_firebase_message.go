package helper

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	cron_service "classroom_itats_api/services/cron"
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type sendFirebaseMessage struct {
	app                  *firebase.App
	presenceCronService  cron_service.PresenceCronService
	taskCronService      cron_service.TaskCronService
	notificationRepo     repositories.NotificationRepository
}

type SendFirebaseMessage interface {
	SendPresenceCreatedNotification() error
	SendPresenceReminderNotification() error
	SendAssignmentCreatedNotification() error
	SendAssignmentReminderNotification() error
}

func NewSendFirebaseMessage(
	app *firebase.App,
	presenceCronService cron_service.PresenceCronService,
	taskCronService cron_service.TaskCronService,
	notificationRepo repositories.NotificationRepository,
) *sendFirebaseMessage {
	return &sendFirebaseMessage{
		app:                 app,
		presenceCronService: presenceCronService,
		taskCronService:     taskCronService,
		notificationRepo:    notificationRepo,
	}
}

// saveNotifForUsers persists a notification row for every recipient.
func (s *sendFirebaseMessage) saveNotifForUsers(ctx context.Context, recipientIDs []string, role, title, body, notifType, referenceID string) {
	for _, id := range recipientIDs {
		n := &entities.Notification{
			RecipientID:   id,
			RecipientRole: role,
			Title:         title,
			Body:          body,
			Type:          notifType,
			ReferenceID:   referenceID,
			IsRead:        false,
			CreatedAt:     time.Now(),
		}
		_ = s.notificationRepo.Save(ctx, n)
	}
}

func (s *sendFirebaseMessage) SendPresenceCreatedNotification() error {
	client, err := s.app.Messaging(context.Background())
	if err != nil {
		return err
	}

	users, err := s.presenceCronService.PresenceCreated(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range users {
		kul := usr["kul"].(entities.Lecture)
		user := usr["user"].([]string)

		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_presence_created_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID))
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			title := "Absensi Baru Telah Dibuat"
			body := fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s Telah dibuat, Silahkan melakukan absensi sebelum tanggal %s", usr["subject"].(string), kul.SubjectClass, kul.PresenceLimit.UTC().Format(time.RFC1123))

			message := &messaging.Message{
				Notification: &messaging.Notification{Title: title, Body: body},
				Topic:        fmt.Sprintf("%s_%s_%s_%s_presence_created_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID),
				Data: map[string]string{
					"type":         "presence",
					"reference_id": kul.LectureID,
				},
			}

			_, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}

			// Persist notification history for each recipient
			s.saveNotifForUsers(context.Background(), user, "student", title, body, "presence", kul.LectureID)
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendPresenceReminderNotification() error {
	client, err := s.app.Messaging(context.Background())
	if err != nil {
		return err
	}

	users, err := s.presenceCronService.PresenceReminder(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range users {
		kul := usr["kul"].(entities.Lecture)
		user := usr["user"].([]string)

		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_presence_reminder_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID))
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			title := "Reminder Absensi"
			body := fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s mendekati batas waktu absensi, Silahkan melakukan absensi sebelum tanggal %s", usr["subject"].(string), kul.SubjectClass, kul.PresenceLimit.UTC().Format(time.RFC1123))

			message := &messaging.Message{
				Notification: &messaging.Notification{Title: title, Body: body},
				Topic:        fmt.Sprintf("%s_%s_%s_%s_presence_reminder_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID),
				Data: map[string]string{
					"type":         "presence",
					"reference_id": kul.LectureID,
				},
			}

			response, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}
			fmt.Println(response, "message were sended successfully")

			s.saveNotifForUsers(context.Background(), user, "student", title, body, "presence", kul.LectureID)
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendAssignmentCreatedNotification() error {
	client, err := s.app.Messaging(context.Background())
	if err != nil {
		return err
	}

	users, err := s.taskCronService.AssignmentCreated(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range users {
		kul := usr["tugaskul"].(entities.Assignment)
		klstw := usr["klstw"].(entities.ClassOffered)
		user := usr["user"].([]string)

		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID))
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			fmt.Println(usr["subject"].(string), kul.SubjectClass, usr, kul, klstw)
			title := "Tugas Baru Telah Dibuat"
			body := fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s Telah dibuat, Silahkan submit tugas sebelum tanggal %s", klstw.SubjectName, klstw.SubjectClass, kul.DueDate.UTC().Format(time.RFC1123))

			message := &messaging.Message{
				Notification: &messaging.Notification{Title: title, Body: body},
				Topic:        fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID),
				Data: map[string]string{
					"type":         "assignment",
					"reference_id": fmt.Sprintf("%d", kul.AssignmentID),
				},
			}

			_, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}

			s.saveNotifForUsers(context.Background(), user, "student", title, body, "assignment", fmt.Sprintf("%d", kul.AssignmentID))
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendAssignmentReminderNotification() error {
	client, err := s.app.Messaging(context.Background())
	if err != nil {
		return err
	}

	users, err := s.taskCronService.AssignmentReminder(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range users {
		kul := usr["tugaskul"].(entities.Assignment)
		klstw := usr["klstw"].(entities.ClassOffered)
		user := usr["user"].([]string)

		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID))
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			title := "Reminder Tugas"
			body := fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s mendekati tenggat pengumpulan, Silahkan submit tugas sebelum tanggal %s", klstw.SubjectName, klstw.SubjectClass, kul.DueDate.UTC().Format(time.RFC1123))

			message := &messaging.Message{
				Notification: &messaging.Notification{Title: title, Body: body},
				Topic:        fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID),
				Data: map[string]string{
					"type":         "assignment",
					"reference_id": fmt.Sprintf("%d", kul.AssignmentID),
				},
			}

			response, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}
			fmt.Println(response, "message were sended successfully")

			s.saveNotifForUsers(context.Background(), user, "student", title, body, "assignment", fmt.Sprintf("%d", kul.AssignmentID))
		}
	}

	return err
}
