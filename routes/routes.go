package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/handlers"
	"pkl-attendance/backend/middleware"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/services"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	studentRepo := repositories.NewStudentRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	journalRepo := repositories.NewJournalRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)
	holidayRepo := repositories.NewHolidayRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	activityRepo := repositories.NewActivityLogRepository(db)

	activitySvc := services.NewActivityLogService(activityRepo)
	authSvc := services.NewAuthService(cfg, userRepo, tokenRepo, activitySvc)
	settingSvc := services.NewSettingService(settingRepo)
	geocoderSvc := services.NewGeocodingService()

	attendanceSvc := services.NewAttendanceService(
		attendanceRepo, studentRepo, journalRepo, leaveRepo, holidayRepo, settingSvc, geocoderSvc,
	)
	journalSvc := services.NewJournalService(journalRepo, attendanceRepo)
	leaveSvc := services.NewLeaveService(leaveRepo)
	studentAdminSvc := services.NewStudentManagementService(userRepo, studentRepo, attendanceRepo, activitySvc)
	adminAttendanceSvc := services.NewAdminAttendanceService(attendanceRepo, activitySvc)
	adminJournalSvc := services.NewAdminJournalService(journalRepo, activitySvc)
	adminLeaveSvc := services.NewAdminLeaveService(leaveRepo, attendanceRepo, holidayRepo, settingSvc, activitySvc)
	holidaySvc := services.NewHolidayService(holidayRepo, activitySvc)
	settingsAdminSvc := services.NewSettingsAdminService(settingSvc, activitySvc)
	recapAdminSvc := services.NewRecapAdminService(attendanceSvc, userRepo)
	dashboardAdminSvc := services.NewDashboardAdminService(attendanceRepo, journalRepo, leaveRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	healthHandler := handlers.NewHealthHandler()
	attendanceHandler := handlers.NewAttendanceHandler(attendanceSvc)
	journalHandler := handlers.NewJournalHandler(journalSvc)
	leaveHandler := handlers.NewLeaveHandler(leaveSvc)
	dashboardHandler := handlers.NewDashboardHandler(attendanceSvc, journalSvc, leaveSvc, settingSvc, userRepo, studentRepo)

	studentAdminHandler := handlers.NewStudentAdminHandler(studentAdminSvc)
	attendanceAdminHandler := handlers.NewAttendanceAdminHandler(adminAttendanceSvc)
	journalAdminHandler := handlers.NewJournalAdminHandler(adminJournalSvc)
	leaveAdminHandler := handlers.NewLeaveAdminHandler(adminLeaveSvc)
	holidayAdminHandler := handlers.NewHolidayAdminHandler(holidaySvc)
	settingsAdminHandler := handlers.NewSettingsAdminHandler(settingsAdminSvc)
	recapAdminHandler := handlers.NewRecapAdminHandler(recapAdminSvc)
	dashboardAdminHandler := handlers.NewDashboardAdminHandler(dashboardAdminSvc)

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.Health)
		loginLimiter := middleware.NewRateLimiter(10, time.Minute)
		api.POST("/auth/login", loginLimiter.Middleware(), authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg, tokenRepo))
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/auth/me", authHandler.Me)
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			siswa := protected.Group("/student")
			siswa.Use(middleware.RequireRole("siswa"))
			{
				siswa.GET("/dashboard", dashboardHandler.SiswaDashboard)
				siswa.POST("/attendance/check-in", attendanceHandler.CheckIn)
				siswa.POST("/attendance/check-out", attendanceHandler.CheckOut)
				siswa.GET("/attendance/today", attendanceHandler.Today)
				siswa.GET("/attendance/history", attendanceHandler.History)
				siswa.GET("/attendance/recap", attendanceHandler.Recap)
				siswa.GET("/journals", journalHandler.ListMine)
				siswa.POST("/journals", journalHandler.Create)
				siswa.PUT("/journals/:id", journalHandler.UpdateMine)
				siswa.GET("/leaves", leaveHandler.ListMine)
				siswa.POST("/leaves", leaveHandler.Create)
			}

			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/dashboard", dashboardAdminHandler.Summary)

				admin.GET("/students", studentAdminHandler.List)
				admin.POST("/students", studentAdminHandler.Create)
				admin.PUT("/students/:id", studentAdminHandler.Update)
				admin.POST("/students/:id/reset-password", studentAdminHandler.ResetPassword)
				admin.PATCH("/students/:id/toggle-active", studentAdminHandler.ToggleActive)
				admin.DELETE("/students/:id", studentAdminHandler.Delete) // <- TAMBAHAN

				admin.GET("/attendances", attendanceAdminHandler.List)
				admin.PUT("/attendances/:id", attendanceAdminHandler.Correct)

				admin.GET("/journals", journalAdminHandler.List)
				admin.POST("/journals/:id/note", journalAdminHandler.AddNote)
				admin.PUT("/journals/:id/override", journalAdminHandler.Override)

				admin.GET("/leaves", leaveAdminHandler.List)
				admin.PUT("/leaves/:id/decide", leaveAdminHandler.Decide)

				admin.GET("/holidays", holidayAdminHandler.List)
				admin.POST("/holidays", holidayAdminHandler.Create)
				admin.PUT("/holidays/:id", holidayAdminHandler.Update)
				admin.DELETE("/holidays/:id", holidayAdminHandler.Delete)

				admin.GET("/settings", settingsAdminHandler.Get)
				admin.PUT("/settings", settingsAdminHandler.Update)

				admin.GET("/recap", recapAdminHandler.RecapAll)
			}
		}
	}
	return r
}