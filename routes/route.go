package routes

import (
	"classroom_itats_api/api/handlers"
	lecturer_handlers "classroom_itats_api/api/handlers/lecturer"
	student_handlers "classroom_itats_api/api/handlers/student"
	"classroom_itats_api/api/middlewares"

	"github.com/gin-gonic/gin"
)

type route struct {
	route                        *gin.Engine
	userHandler                  handlers.UserHandler
	studentSubjectHandler        student_handlers.StudentSubjectHandler
	lecturerSubjectHandler       lecturer_handlers.LecturerSubjectHandler
	forumHandler                 handlers.ForumHandler
	lectureHandler               handlers.LectureHandler
	studentPresenceHandler       student_handlers.StudentPresenceHandler
	studentMaterialHandler       student_handlers.StudentMaterialHandler
	subjectMemberHandler         handlers.SubjectMemberHandler
	studentProfileHandler        student_handlers.StudentProfileHandler
	lecturerAssignmentHandler    lecturer_handlers.LecturerAssignmentHandler
	lecturerCollegeReportHandler lecturer_handlers.LecturerCollegeReportHandler
	lecturerMaterialHandler      lecturer_handlers.LecturerMaterialHandler
	webhookHandler               handlers.WebhookHandler
}

type Route interface {
	Routes()
}

func NewRoute(run *gin.Engine,
	userHandler handlers.UserHandler,
	studentSubjectHandler student_handlers.StudentSubjectHandler,
	lecturerSubjectHandler lecturer_handlers.LecturerSubjectHandler,
	forumHandler handlers.ForumHandler, lectureHandler handlers.LectureHandler,
	studentPresenceHandler student_handlers.StudentPresenceHandler,
	studentMaterialHandler student_handlers.StudentMaterialHandler,
	subjectMemberHandler handlers.SubjectMemberHandler,
	studentProfileHandler student_handlers.StudentProfileHandler,
	lecturerAssignmentHandler lecturer_handlers.LecturerAssignmentHandler,
	lecturerCollegeReportHandler lecturer_handlers.LecturerCollegeReportHandler,
	lecturerMaterialHandler lecturer_handlers.LecturerMaterialHandler,
	webhookHandler handlers.WebhookHandler,
) *route {
	return &route{
		route:                        run,
		userHandler:                  userHandler,
		studentSubjectHandler:        studentSubjectHandler,
		lecturerSubjectHandler:       lecturerSubjectHandler,
		forumHandler:                 forumHandler,
		lectureHandler:               lectureHandler,
		studentPresenceHandler:       studentPresenceHandler,
		studentMaterialHandler:       studentMaterialHandler,
		subjectMemberHandler:         subjectMemberHandler,
		studentProfileHandler:        studentProfileHandler,
		lecturerAssignmentHandler:    lecturerAssignmentHandler,
		lecturerCollegeReportHandler: lecturerCollegeReportHandler,
		lecturerMaterialHandler:      lecturerMaterialHandler,
		webhookHandler:               webhookHandler,
	}
}

func (r *route) Routes() {
	public := r.route.Group("/api/v1")
	public.POST("/login", r.userHandler.Login)
	public.PUT("/login/info", r.userHandler.StoreLoginInfo)
	public.PUT("/logout", r.userHandler.Logout)

	verify := public.Group("/verify-npm")
	verify.Use(middlewares.ApiKey())
	verify.POST("/", r.userHandler.CheckNPMIsExist)
	verifys := public.Group("/verify-npms")
	verifys.Use(middlewares.ApiKey())
	verifys.POST("/", r.userHandler.CheckNPMsIsExist)

	private := r.route.Group("/api/v1")
	private.Use(middlewares.Auth())
	private.POST("/all/student", r.userHandler.GetDataMhs)
	private.POST("/subjects/forums", r.forumHandler.Forum)
	private.POST("/subjects/forums/store", r.forumHandler.StoreAnnouncement)
	private.PUT("/subjects/forums/update", r.forumHandler.UpdateAnnouncement)
	private.DELETE("/subjects/forums/delete", r.forumHandler.DeleteAnnouncement)
	private.POST("/subjects/forums/comments/store", r.forumHandler.StoreComment)
	private.PUT("/subjects/forums/comments/update", r.forumHandler.UpdateComment)
	private.DELETE("/subjects/forums/comments/delete", r.forumHandler.DeleteComment)
	private.POST("/subjects/lectures", r.lectureHandler.Getlecture)
	private.POST("/subjects/members", r.subjectMemberHandler.SubjectMember)

	student := private.Group("/students")
	student.GET("/subjects", r.studentSubjectHandler.GetActiveSubjectByStudentWithID)
	student.POST("/subjects", r.studentSubjectHandler.GetSubjectByStudentIdWithPeriod)
	student.GET("/periodes", r.studentSubjectHandler.StudyPeriodes)
	student.POST("/profile", r.studentProfileHandler.GetStudentProfile)
	student.GET("/home/presences/active", r.studentPresenceHandler.GetHomeActivePresence)
	student.GET("/home/assignments/active", r.studentMaterialHandler.GetHomeActiveAssignment)
	subjects := student.Group("/subjects")
	subjects.POST("/lectures", r.lectureHandler.GetStudentLecture)
	subjects.POST("/responsi", r.studentPresenceHandler.GetSubjectResponsi)
	subjects.POST("/presences", r.studentPresenceHandler.GetStudentPresences)
	subjects.POST("/presences/present", r.studentPresenceHandler.SetPresenceStudent)
	subjects.POST("/presences/active", r.studentPresenceHandler.GetActivePresence)
	subjects.POST("/presences/questions", r.studentPresenceHandler.GetPresenceQuestion)
	subjects.POST("/materials", r.studentMaterialHandler.GetStudentMaterial)
	subjects.POST("/materials/weeks", r.studentMaterialHandler.GetWeekMaterial)
	subjects.POST("/materials/achievements", r.studentMaterialHandler.GetStudyAchievement)
	subjects.POST("/materials/assignments", r.studentMaterialHandler.GetStudentAssignmentGroup)
	subjects.POST("/materials/assignments/detail", r.studentMaterialHandler.GetStudentAssignment)
	subjects.POST("/materials/assignments/submited", r.studentMaterialHandler.GetStudentAssignmentSubmission)
	subjects.POST("/materials/assignments/active", r.studentMaterialHandler.GetActiveAssignment)
	subjects.POST("/scores", r.studentMaterialHandler.GetStudentAssignmentScores)

	lecturer := private.Group("/lecturers")
	lecturer.GET("/subjects", r.lecturerSubjectHandler.GetActiveSubjectByLecturerID)
	lecturer.POST("/subjects", r.lecturerSubjectHandler.GetSubjectByLecturerIDFiltered)
	lecturer.GET("/periodes", r.lecturerSubjectHandler.LecturePeriodes)
	lecturer.POST("/majors", r.lecturerSubjectHandler.LecturerSubjectMajor)
	lecturer.POST("/assignments", r.lecturerAssignmentHandler.GetLecturerCreatedAssignment)
	lecturer.GET("/assignments/weeks", r.lecturerAssignmentHandler.GetWeekAssignment)
	lecturer.GET("/assignments/scoreType", r.lecturerAssignmentHandler.GetScoreTypeAssignment)
	lecturer.POST("/studentScores", r.lecturerSubjectHandler.GetStudentScore)
	lecturer.POST("/percentages", r.lecturerSubjectHandler.GetPercentage)
	lecturer.GET("/subjects/reports", r.lecturerSubjectHandler.GetActiveSubjectReportByLecturerID)
	lecturer.POST("/subjects/reports", r.lecturerSubjectHandler.GetSubjectReportByLecturerIDFiltered)
	lecturer.POST("/colleges/reports", r.lecturerCollegeReportHandler.GetSubjectCollegeReport)
	lecturer.POST("/colleges/reports/store", r.lecturerCollegeReportHandler.CreateCollege)
	lecturer.PUT("/colleges/reports/edit", r.lecturerCollegeReportHandler.EditCollege)
	lecturer.DELETE("/colleges/reports/delete", r.lecturerCollegeReportHandler.DeleteCollege)
	lecturer.POST("/colleges/reports/detail", r.lecturerCollegeReportHandler.GetSubjectCollegeReportByKulID)
	lecturer.POST("/colleges/materials", r.lecturerMaterialHandler.GetMaterialSelected)
	lecturer.GET("/materials", r.lecturerMaterialHandler.GetMaterials)
	lecturer.POST("/colleges/team-weeks", r.lecturerCollegeReportHandler.GetTeamWeeks)
	lecturer.GET("/colleges/rps", r.lecturerCollegeReportHandler.GetRPSDetail)

	// Webhook endpoints — diamankan dengan x_api_key
	webhook := r.route.Group("/api/v1/internal")
	webhook.Use(middlewares.ApiKey())
	webhook.POST("/notify/presence", r.webhookHandler.TriggerPresenceNotification)
	webhook.POST("/notify/assignment", r.webhookHandler.TriggerAssignmentNotification)
}
