package helper

import (
	"classroom_itats_api/entities"
	cron_service "classroom_itats_api/services/cron"
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type sendFirebaseMessage struct {
	app                 *firebase.App
	presenceCronService cron_service.PresenceCronService
	taskCronService     cron_service.TaskCronService
}

type SendFirebaseMessage interface {
	SendPresenceCreatedNotification() error
	SendPresenceReminderNotification() error
	SendAssignmentCreatedNotification() error
	SendAssignmentReminderNotification() error
}

func NewSendFirebaseMessage(app *firebase.App, presenceCronService cron_service.PresenceCronService, taskCronService cron_service.TaskCronService) *sendFirebaseMessage {
	return &sendFirebaseMessage{app: app, presenceCronService: presenceCronService, taskCronService: taskCronService}
}

func (s *sendFirebaseMessage) SendPresenceCreatedNotification() error {
	// Obtain a messaging client from the Firebase app
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

		// These registration tokens come from the client FCM SDKs.
		user := usr["user"].([]string)

		// Subscribe the devices corresponding to the registration tokens to the
		// topic.
		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_presence_created_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID))

		// See the TopicManagementResponse reference documentation
		// for the contents of response.
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")
		if e != nil {
			err = e
		} else {
			// See documentation on defining a message payload.
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "Absensi Baru Telah Dibuat",
					Body:  fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s Telah dibuat, Silahkan melakukan absensi sebelum tanggal %s", usr["subject"].(string), kul.SubjectClass, kul.PresenceLimit.UTC().Format(time.RFC1123)),
				},
				Topic: fmt.Sprintf("%s_%s_%s_%s_presence_created_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID),
			}

			// Send a message to the devices subscribed to the provided topic.
			_, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendPresenceReminderNotification() error {
	// Obtain a messaging client from the Firebase app
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

		// These registration tokens come from the client FCM SDKs.
		user := usr["user"].([]string)

		// Subscribe the devices corresponding to the registration tokens to the
		// topic.
		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_presence_reminder_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID))

		// See the TopicManagementResponse reference documentation
		// for the contents of response.
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			// See documentation on defining a message payload.
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "Reminder Absensi",
					Body:  fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s mendekati batas waktu absensi, Silahkan melakukan absensi sebelum tanggal %s", usr["subject"].(string), kul.SubjectClass, kul.PresenceLimit.UTC().Format(time.RFC1123)),
				},
				Topic: fmt.Sprintf("%s_%s_%s_%s_presence_reminder_notification", kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID),
			}

			// Send a message to the devices subscribed to the provided topic.
			response, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}

			fmt.Println(response, "message were sended successfully")
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendAssignmentCreatedNotification() error {
	// Obtain a messaging client from the Firebase app
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

		// These registration tokens come from the client FCM SDKs.
		user := usr["user"].([]string)

		// Subscribe the devices corresponding to the registration tokens to the
		// topic.
		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID))

		// See the TopicManagementResponse reference documentation
		// for the contents of response.
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")
		if e != nil {
			err = e
		} else {
			// See documentation on defining a message payload.
			fmt.Println(usr["subject"].(string), kul.SubjectClass, usr, kul, klstw)
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "Tugas Baru Telah Dibuat",
					Body:  fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s Telah dibuat, Silahkan submit tugas sebelum tanggal %s", klstw.SubjectName, klstw.SubjectClass, kul.DueDate.UTC().Format(time.RFC1123)),
				},
				Topic: fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID),
			}

			// Send a message to the devices subscribed to the provided topic.
			_, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}
		}
	}

	return err
}

func (s *sendFirebaseMessage) SendAssignmentReminderNotification() error {
	// Obtain a messaging client from the Firebase app
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

		// These registration tokens come from the client FCM SDKs.
		user := usr["user"].([]string)

		// Subscribe the devices corresponding to the registration tokens to the
		// topic.
		response, e := client.SubscribeToTopic(context.Background(), user, fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID))

		// See the TopicManagementResponse reference documentation
		// for the contents of response.
		fmt.Println(response.SuccessCount, "tokens were subscribed successfully")

		if e != nil {
			err = e
		} else {
			// See documentation on defining a message payload.
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: "Reminder Tugas",
					Body:  fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s mendekati tenggat pengumpulan, Silahkan submit tugas sebelum tanggal %s", klstw.SubjectName, klstw.SubjectClass, kul.DueDate.UTC().Format(time.RFC1123)),
				},
				Topic: fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification", klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID),
			}

			// Send a message to the devices subscribed to the provided topic.
			response, e := client.Send(context.Background(), message)
			if e != nil {
				err = e
			}

			fmt.Println(response, "message were sended successfully")
		}
	}

	return err
}
