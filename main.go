package main

import (
	"classroom_itats_api/api/handlers"
	lecturer_handlers "classroom_itats_api/api/handlers/lecturer"
	student_handlers "classroom_itats_api/api/handlers/student"
	"classroom_itats_api/helper"
	"classroom_itats_api/repositories"
	lecturer_repositories "classroom_itats_api/repositories/lecturer"
	student_repositories "classroom_itats_api/repositories/student"
	"classroom_itats_api/routes"
	"classroom_itats_api/services"
	cron_service "classroom_itats_api/services/cron"
	lecturer_services "classroom_itats_api/services/lecturer"
	microsoft_services "classroom_itats_api/services/microsoft"
	student_services "classroom_itats_api/services/student"
	"net/http"
	"time"

	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ==================================================================
// https://firebase.google.com/docs/admin/setup
// ==================================================================

func initializeServiceAccountID() *firebase.App {
	// [START initialize_sdk_with_service_account_id]
	conf := &firebase.Config{
		ServiceAccountID: "firebase-adminsdk-fowjd@871324361748.iam.gserviceaccount.com",
	}
	opt := option.WithCredentialsFile("./resource/classroomitats-firebase-adminsdk-fowjd-1098de7885.json")
	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_sdk_with_service_account_id]
	return app
}

func main() {
	loadEnv()

	run := gin.Default()

	db := helper.NewDatabaseConnection()
	conn, _ := db.Connect()

	userRepository := repositories.NewUserRepository(conn)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	studentSubjectRepository := student_repositories.NewStudentSubjectRepository(conn)
	studentSubjectService := student_services.NewStudentSubjectService(studentSubjectRepository)
	studentSubjectHandler := student_handlers.NewStudentSubjectService(studentSubjectService)

	lecturerSubjectRepository := lecturer_repositories.NewLecturerSubjectRepository(conn)
	lecturerSubjectService := lecturer_services.NewLecturerSubjectService(lecturerSubjectRepository)
	lecturerSubjecthandler := lecturer_handlers.NewLecturerSubjectHanlder(lecturerSubjectService)

	forumRepository := repositories.NewForumRepository(conn)
	forumService := services.NewForumService(forumRepository)
	forumhandler := handlers.NewForumHandler(forumService)

	lectureRepository := repositories.NewLectureRepository(conn)
	lectureService := services.NewLectureService(lectureRepository)
	lecturehandler := handlers.NewLectureHandler(lectureService)

	studentPresenceRepository := student_repositories.NewStudentPresenceRepository(conn)
	studentPresenceService := student_services.NewStudentPresenceService(studentPresenceRepository)
	studentPresenceHandler := student_handlers.NewStudentPresenceHandler(studentPresenceService)

	studentMaterialRepository := student_repositories.NewStudentMaterialRepository(conn)
	studentMaterialService := student_services.NewStudentMaterialService(studentMaterialRepository)
	studentMaterialHandler := student_handlers.NewStudentMaterialHandler(studentMaterialService)

	subjectMemberRepository := repositories.NewSubjectMemberRepository(conn)
	subjectMemberService := services.NewSubjectMemberService(subjectMemberRepository)
	subjectMemberHandler := handlers.NewSubjectMemberHandler(subjectMemberService)

	studentProfileRepository := student_repositories.NewStudentProfileRepository(conn)
	studentProfileService := student_services.NewStudentProfileService(studentProfileRepository)
	studentProfileHandler := student_handlers.NewStudentProfileHandler(studentProfileService)

	lecturerAssignmentRepository := lecturer_repositories.NewLecturerAssignmentRepository(conn)
	lecturerAssignmentService := lecturer_services.NewLecturerAssignmentService(lecturerAssignmentRepository)
	lecturerAssignmenthandler := lecturer_handlers.NewLecturerAssignmentHandlder(lecturerAssignmentService)

	lecturerCollegeReportRepository := lecturer_repositories.NewLecturerCollegeReportRepository(conn)
	lecturerCollegeReportService := lecturer_services.NewLecturerCollegeReportService(lecturerCollegeReportRepository)
	lecturerCollegeReporthandler := lecturer_handlers.NewLecturerCollegeReportHandlder(lecturerCollegeReportService)

	lecturerMaterialRepository := lecturer_repositories.NewLecturerMaterialRepository(conn)
	lecturerMaterialService := lecturer_services.NewLecturerMaterialService(lecturerMaterialRepository)
	lecturerMaterialhandler := lecturer_handlers.NewLecturerMaterialHandlder(lecturerMaterialService)

	// Microsoft Teams OAuth Service
	msAuthService := microsoft_services.NewMicrosoftAuthService(conn)
	lecturerMicrosoftHandler := lecturer_handlers.NewLecturerMicrosoftHandler(msAuthService)

	firebaseApp := initializeServiceAccountID()

	presenceCronService := cron_service.NewPresenceCronService(studentPresenceRepository)
	taskCronService := cron_service.NewTaskCronService(studentMaterialRepository)

	notificationRepo := repositories.NewNotificationRepository(conn)
	notificationHandler := handlers.NewNotificationHandler(notificationRepo)

	firebaseHelper := helper.NewSendFirebaseMessage(firebaseApp, presenceCronService, taskCronService, notificationRepo)

	webhookHandler := handlers.NewWebhookHandler(firebaseHelper)

	route := routes.NewRoute(run,
		userHandler,
		studentSubjectHandler,
		lecturerSubjecthandler,
		forumhandler,
		lecturehandler,
		studentPresenceHandler,
		studentMaterialHandler,
		subjectMemberHandler,
		studentProfileHandler,
		lecturerAssignmenthandler,
		lecturerCollegeReporthandler,
		lecturerMaterialhandler,
		lecturerMicrosoftHandler,
		webhookHandler,
		notificationHandler,
	)

	route.Routes()

	timezone, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		panic(err.Error())
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(timezone))
	if err != nil {
		panic(err.Error())
	}

	j1, err := s.NewJob(gocron.DurationJob(time.Minute*30), gocron.NewTask(func() {
		firebaseHelper.SendPresenceCreatedNotification()
		firebaseHelper.SendAssignmentCreatedNotification()
	}))

	if err != nil {
		panic(err.Error())
	}

	j2, err := s.NewJob(gocron.DurationJob(time.Hour*3), gocron.NewTask(func() {
		firebaseHelper.SendPresenceReminderNotification()
		firebaseHelper.SendAssignmentReminderNotification()
	}))

	if err != nil {
		panic(err.Error())
	}

	log.Println(j1.LastRun())
	log.Println(j2.LastRun())

	s.Start()

	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	// err = s.Shutdown()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// run.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://127.0.0.1"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "http://127.0.0.1:/"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))

	server := &http.Server{
		Addr:           viper.GetString("URL"),
		Handler:        run,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 2048,
	}

	server.ListenAndServe()
}

func loadEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}
}