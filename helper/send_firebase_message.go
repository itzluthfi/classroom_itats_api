package helper

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	cron_service "classroom_itats_api/services/cron"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type sendFirebaseMessage struct {
	app                 *firebase.App
	presenceCronService cron_service.PresenceCronService
	taskCronService     cron_service.TaskCronService
	notificationRepo    repositories.NotificationRepository
}

type SendFirebaseMessage interface {
	SendPresenceCreatedNotification() error
	SendPresenceReminderNotification() error
	SendAssignmentCreatedNotification() error
	SendAssignmentReminderNotification() error
	SendAssignmentReminderH1Notification() error
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

// saveNotifForUsers menyimpan notifikasi ke DB untuk setiap penerima.
// Error discard karena tidak boleh memblok alur utama pengiriman FCM.
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
		if saveErr := s.notificationRepo.Save(ctx, n); saveErr != nil {
			log.Printf("[NOTIF] gagal simpan notifikasi untuk %s: %v", id, saveErr)
		}
	}
}

// extractLecture mengambil entities.Lecture dari map dengan type-assertion aman.
func extractLecture(usr map[string]interface{}) (entities.Lecture, error) {
	raw, ok := usr["kul"]
	if !ok {
		return entities.Lecture{}, errors.New("key 'kul' tidak ditemukan di data cron presence")
	}
	kul, ok := raw.(entities.Lecture)
	if !ok {
		return entities.Lecture{}, fmt.Errorf("type assertion gagal untuk 'kul': got %T", raw)
	}
	return kul, nil
}

// extractStringSlice mengambil []string dari map dengan type-assertion aman.
func extractStringSlice(usr map[string]interface{}, key string) ([]string, error) {
	raw, ok := usr[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' tidak ditemukan", key)
	}
	val, ok := raw.([]string)
	if !ok {
		return nil, fmt.Errorf("type assertion gagal untuk '%s': got %T", key, raw)
	}
	return val, nil
}

// extractString mengambil string dari map dengan type-assertion aman.
func extractString(usr map[string]interface{}, key string) string {
	raw, ok := usr[key]
	if !ok {
		return ""
	}
	val, ok := raw.(string)
	if !ok {
		return ""
	}
	return val
}

// formatTime memformat time.Time dengan aman — menghindari zero value.
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "waktu tidak tersedia"
	}
	return t.UTC().Format(time.RFC1123)
}

// ─── SendPresenceCreatedNotification ─────────────────────────────────────────

func (s *sendFirebaseMessage) SendPresenceCreatedNotification() error {
	ctx := context.Background()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		log.Printf("[FCM Error] Gagal inisialisasi client: %v\n", err)
		return fmt.Errorf("gagal inisialisasi Firebase Messaging client: %w", err)
	}

	users, err := s.presenceCronService.PresenceCreated(ctx)
	if err != nil {
		log.Printf("[FCM Error] Gagal ambil data presence dari DB: %v\n", err)
		return fmt.Errorf("gagal ambil data presence created dari DB: %w", err)
	}

	if len(users) == 0 {
		log.Println("[FCM] tidak ada absensi baru yang perlu dinotifikasi")
		return nil
	}

	var errs []error

	for _, usr := range users {
		kul, e := extractLecture(usr)
		if e != nil {
			errs = append(errs, fmt.Errorf("[PresenceCreated] data tidak valid: %w", e))
			continue
		}

		user, e := extractStringSlice(usr, "user")
		if e != nil {
			errs = append(errs, fmt.Errorf("[PresenceCreated] kul=%s: %w", kul.LectureID, e))
			continue
		}

		// Skip jika tidak ada token — Firebase error jika slice kosong
		if len(user) == 0 {
			log.Printf("[FCM] PresenceCreated: tidak ada mahasiswa dengan token untuk kulid=%s", kul.LectureID)
			continue
		}

		topic := fmt.Sprintf("%s_%s_%s_%s_presence_created_notification",
			kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID)

		log.Printf("[FCM] Memproses presence untuk topic: %s dengan %d token", topic, len(user))

		response, e := client.SubscribeToTopic(ctx, user, topic)
		if e != nil {
			log.Printf("[FCM Error] SubscribeToTopic gagal: %v\n", e)
			errs = append(errs, fmt.Errorf("[PresenceCreated] SubscribeToTopic gagal untuk kulid=%s: %w", kul.LectureID, e))
			continue
		}
		log.Printf("[FCM] PresenceCreated: %d token berhasil subscribe ke topic %s", response.SuccessCount, topic)

		subject := extractString(usr, "subject")
		if subject == "" {
			subject = "Matakuliah"
		}
		title := "Absensi Baru Telah Dibuat"
		body := fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s telah dibuat. Silahkan absensi sebelum %s",
			subject, kul.SubjectClass, formatTime(kul.PresenceLimit))

		message := &messaging.Message{
			Notification: &messaging.Notification{Title: title, Body: body},
			Topic:        topic,
			Data: map[string]string{
				"type":         "presence",
				"reference_id": kul.LectureID,
			},
		}

		if _, e = client.Send(ctx, message); e != nil {
			errs = append(errs, fmt.Errorf("[PresenceCreated] Send FCM gagal untuk kulid=%s: %w", kul.LectureID, e))
			continue
		}

		s.saveNotifForUsers(ctx, user, "student", title, body, "presence", kul.LectureID)
	}

	return errors.Join(errs...)
}

// ─── SendPresenceReminderNotification ────────────────────────────────────────

func (s *sendFirebaseMessage) SendPresenceReminderNotification() error {
	ctx := context.Background()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("gagal inisialisasi Firebase Messaging client: %w", err)
	}

	users, err := s.presenceCronService.PresenceReminder(ctx)
	if err != nil {
		return fmt.Errorf("gagal ambil data presence reminder dari DB: %w", err)
	}

	if len(users) == 0 {
		log.Println("[FCM] tidak ada absensi yang mendekati batas waktu")
		return nil
	}

	var errs []error

	for _, usr := range users {
		kul, e := extractLecture(usr)
		if e != nil {
			errs = append(errs, fmt.Errorf("[PresenceReminder] data tidak valid: %w", e))
			continue
		}

		user, e := extractStringSlice(usr, "user")
		if e != nil {
			errs = append(errs, fmt.Errorf("[PresenceReminder] kul=%s: %w", kul.LectureID, e))
			continue
		}

		if len(user) == 0 {
			log.Printf("[FCM] PresenceReminder: tidak ada mahasiswa dengan token untuk kulid=%s", kul.LectureID)
			continue
		}

		topic := fmt.Sprintf("%s_%s_%s_%s_presence_reminder_notification",
			kul.SubjectID, kul.SubjectClass, kul.AcademicPeriodID, kul.MajorID)

		response, e := client.SubscribeToTopic(ctx, user, topic)
		if e != nil {
			errs = append(errs, fmt.Errorf("[PresenceReminder] SubscribeToTopic gagal untuk kulid=%s: %w", kul.LectureID, e))
			continue
		}
		log.Printf("[FCM] PresenceReminder: %d token subscribe ke topic %s", response.SuccessCount, topic)

		subject := extractString(usr, "subject")
		title := "Reminder Absensi"
		body := fmt.Sprintf("Absensi Mata Kuliah %s Kelas %s mendekati batas waktu. Silahkan absensi sebelum %s",
			subject, kul.SubjectClass, formatTime(kul.PresenceLimit))

		message := &messaging.Message{
			Notification: &messaging.Notification{Title: title, Body: body},
			Topic:        topic,
			Data: map[string]string{
				"type":         "presence",
				"reference_id": kul.LectureID,
			},
		}

		if _, e = client.Send(ctx, message); e != nil {
			errs = append(errs, fmt.Errorf("[PresenceReminder] Send FCM gagal untuk kulid=%s: %w", kul.LectureID, e))
			continue
		}

		s.saveNotifForUsers(ctx, user, "student", title, body, "presence", kul.LectureID)
	}

	return errors.Join(errs...)
}

// ─── SendAssignmentCreatedNotification ───────────────────────────────────────

func (s *sendFirebaseMessage) SendAssignmentCreatedNotification() error {
	ctx := context.Background()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		log.Printf("[FCM Error] Gagal inisialisasi client: %v\n", err)
		return fmt.Errorf("gagal inisialisasi Firebase Messaging client: %w", err)
	}

	users, err := s.taskCronService.AssignmentCreated(ctx)
	if err != nil {
		log.Printf("[FCM Error] Gagal ambil data assignment dari DB: %v\n", err)
		return fmt.Errorf("gagal ambil data assignment created dari DB: %w", err)
	}

	if len(users) == 0 {
		log.Println("[FCM] tidak ada tugas baru yang perlu dinotifikasi")
		return nil
	}

	var errs []error

	for _, usr := range users {
		// Type assertion aman untuk Assignment
		rawKul, ok := usr["tugaskul"]
		if !ok {
			errs = append(errs, errors.New("[AssignmentCreated] key 'tugaskul' tidak ditemukan"))
			continue
		}
		kul, ok := rawKul.(entities.Assignment)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] type assertion gagal untuk 'tugaskul': got %T", rawKul))
			continue
		}

		// Type assertion aman untuk ClassOffered
		rawKlstw, ok := usr["klstw"]
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] key 'klstw' tidak ditemukan untuk assignment_id=%d", kul.AssignmentID))
			continue
		}
		klstw, ok := rawKlstw.(entities.ClassOffered)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] type assertion gagal untuk 'klstw': got %T", rawKlstw))
			continue
		}

		user, e := extractStringSlice(usr, "user")
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		if len(user) == 0 {
			log.Printf("[FCM] AssignmentCreated: tidak ada mahasiswa dengan token untuk assignment_id=%d", kul.AssignmentID)
			continue
		}

		topic := fmt.Sprintf("%s_%s_%s_%s_%d_assignment_created_notification",
			klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID)

		response, e := client.SubscribeToTopic(ctx, user, topic)
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] SubscribeToTopic gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}
		log.Printf("[FCM] AssignmentCreated: %d token subscribe ke topic %s", response.SuccessCount, topic)

		title := "Tugas Baru Telah Dibuat"
		body := fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s telah dibuat. Silahkan submit sebelum %s",
			klstw.SubjectName, klstw.SubjectClass, formatTime(kul.DueDate))

		message := &messaging.Message{
			Notification: &messaging.Notification{Title: title, Body: body},
			Topic:        topic,
			Data: map[string]string{
				"type":         "assignment",
				"reference_id": fmt.Sprintf("%d", kul.AssignmentID),
			},
		}

		if _, e = client.Send(ctx, message); e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentCreated] Send FCM gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		s.saveNotifForUsers(ctx, user, "student", title, body, "assignment", fmt.Sprintf("%d", kul.AssignmentID))
	}

	return errors.Join(errs...)
}

// ─── SendAssignmentReminderNotification ──────────────────────────────────────

func (s *sendFirebaseMessage) SendAssignmentReminderNotification() error {
	ctx := context.Background()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("gagal inisialisasi Firebase Messaging client: %w", err)
	}

	users, err := s.taskCronService.AssignmentReminder(ctx)
	if err != nil {
		return fmt.Errorf("gagal ambil data assignment reminder dari DB: %w", err)
	}

	if len(users) == 0 {
		log.Println("[FCM] tidak ada tugas yang mendekati tenggat")
		return nil
	}

	var errs []error

	for _, usr := range users {
		rawKul, ok := usr["tugaskul"]
		if !ok {
			errs = append(errs, errors.New("[AssignmentReminder] key 'tugaskul' tidak ditemukan"))
			continue
		}
		kul, ok := rawKul.(entities.Assignment)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] type assertion gagal untuk 'tugaskul': got %T", rawKul))
			continue
		}

		rawKlstw, ok := usr["klstw"]
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] key 'klstw' tidak ditemukan untuk assignment_id=%d", kul.AssignmentID))
			continue
		}
		klstw, ok := rawKlstw.(entities.ClassOffered)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] type assertion gagal untuk 'klstw': got %T", rawKlstw))
			continue
		}

		user, e := extractStringSlice(usr, "user")
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		if len(user) == 0 {
			log.Printf("[FCM] AssignmentReminder: tidak ada mahasiswa dengan token untuk assignment_id=%d", kul.AssignmentID)
			continue
		}

		topic := fmt.Sprintf("%s_%s_%s_%s_%d_assignment_reminder_notification",
			klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID)

		response, e := client.SubscribeToTopic(ctx, user, topic)
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] SubscribeToTopic gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}
		log.Printf("[FCM] AssignmentReminder: %d token subscribe ke topic %s", response.SuccessCount, topic)

		title := "Reminder Tugas"
		body := fmt.Sprintf("Tugas Mata Kuliah %s Kelas %s mendekati tenggat pengumpulan. Silahkan submit sebelum %s",
			klstw.SubjectName, klstw.SubjectClass, formatTime(kul.DueDate))

		message := &messaging.Message{
			Notification: &messaging.Notification{Title: title, Body: body},
			Topic:        topic,
			Data: map[string]string{
				"type":         "assignment",
				"reference_id": fmt.Sprintf("%d", kul.AssignmentID),
			},
		}

		if _, e = client.Send(ctx, message); e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminder] Send FCM gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		s.saveNotifForUsers(ctx, user, "student", title, body, "assignment", fmt.Sprintf("%d", kul.AssignmentID))
	}

	return errors.Join(errs...)
}

// ─── SendAssignmentReminderH1Notification ────────────────────────────────────

func (s *sendFirebaseMessage) SendAssignmentReminderH1Notification() error {
	ctx := context.Background()

	client, err := s.app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("gagal inisialisasi Firebase Messaging client: %w", err)
	}

	users, err := s.taskCronService.AssignmentReminderH1(ctx)
	if err != nil {
		return fmt.Errorf("gagal ambil data assignment reminder H-1 dari DB: %w", err)
	}

	if len(users) == 0 {
		log.Println("[FCM] tidak ada tugas yang mendekati deadline H-1 besok")
		return nil
	}

	var errs []error

	for _, usr := range users {
		rawKul, ok := usr["tugaskul"]
		if !ok {
			errs = append(errs, errors.New("[AssignmentReminderH1] key 'tugaskul' tidak ditemukan"))
			continue
		}
		kul, ok := rawKul.(entities.Assignment)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] type assertion gagal untuk 'tugaskul': got %T", rawKul))
			continue
		}

		rawKlstw, ok := usr["klstw"]
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] key 'klstw' tidak ditemukan untuk assignment_id=%d", kul.AssignmentID))
			continue
		}
		klstw, ok := rawKlstw.(entities.ClassOffered)
		if !ok {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] type assertion gagal untuk 'klstw': got %T", rawKlstw))
			continue
		}

		user, e := extractStringSlice(usr, "user")
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		if len(user) == 0 {
			log.Printf("[FCM] AssignmentReminderH1: tidak ada mahasiswa dengan token untuk assignment_id=%d", kul.AssignmentID)
			continue
		}

		topic := fmt.Sprintf("%s_%s_%s_%s_%d_assignment_reminder_h1_notification",
			klstw.SubjectID, klstw.SubjectClass, klstw.AcademicPeriodID, klstw.MajorID, kul.AssignmentID)

		response, e := client.SubscribeToTopic(ctx, user, topic)
		if e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] SubscribeToTopic gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}
		log.Printf("[FCM] AssignmentReminderH1: %d token subscribe ke topic %s", response.SuccessCount, topic)

		title := "Awas Deadline Besok! 🚨"
		body := fmt.Sprintf("Awas! Tugas %s deadline besok malam (%s) dan kamu belum kumpul.",
			kul.AssignmentTitle, formatTime(kul.DueDate))

		message := &messaging.Message{
			Notification: &messaging.Notification{Title: title, Body: body},
			Topic:        topic,
			Data: map[string]string{
				"type":         "assignment",
				"reference_id": fmt.Sprintf("%d", kul.AssignmentID),
			},
		}

		if _, e = client.Send(ctx, message); e != nil {
			errs = append(errs, fmt.Errorf("[AssignmentReminderH1] Send FCM gagal untuk assignment_id=%d: %w", kul.AssignmentID, e))
			continue
		}

		s.saveNotifForUsers(ctx, user, "student", title, body, "assignment", fmt.Sprintf("%d", kul.AssignmentID))
	}

	return errors.Join(errs...)
}
