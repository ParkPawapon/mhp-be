package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/http/handlers"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type Dependencies struct {
	Config              config.Config
	Logger              *zap.Logger
	DB                  *gorm.DB
	Redis               *redis.Client
	AuthService         services.AuthService
	UserService         services.UserService
	CaregiverService    services.CaregiverService
	MedicineService     services.MedicineService
	IntakeService       services.IntakeService
	HealthService       services.HealthService
	AppointmentService  services.AppointmentService
	ContentService      services.ContentService
	NotificationService services.NotificationService
	SupportService      services.SupportService
	AdminService        services.AdminService
	AuditService        services.AuditService
}

func NewRouter(deps Dependencies) *gin.Engine {
	if deps.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(deps.Logger))
	r.Use(middleware.Recovery(deps.Logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     deps.Config.CORS.OriginsList(),
		AllowMethods:     deps.Config.CORS.MethodsList(),
		AllowHeaders:     deps.Config.CORS.HeadersList(),
		AllowCredentials: deps.Config.CORS.AllowCredentials,
	}))

	healthHandler := handlers.NewHealthHandler(deps.DB, deps.Redis)
	authHandler := handlers.NewAuthHandler(deps.AuthService)
	userHandler := handlers.NewUserHandler(deps.UserService)
	caregiverHandler := handlers.NewCaregiverHandler(deps.CaregiverService)
	medicineHandler := handlers.NewMedicineHandler(deps.MedicineService)
	intakeHandler := handlers.NewIntakeHandler(deps.IntakeService, deps.CaregiverService)
	healthRecordHandler := handlers.NewHealthRecordsHandler(deps.HealthService, deps.CaregiverService)
	appointmentHandler := handlers.NewAppointmentHandler(deps.AppointmentService, deps.CaregiverService)
	contentHandler := handlers.NewContentHandler(deps.ContentService)
	notificationHandler := handlers.NewNotificationHandler(deps.NotificationService)
	supportHandler := handlers.NewSupportHandler(deps.SupportService)
	adminHandler := handlers.NewAdminHandler(deps.AdminService)
	auditHandler := handlers.NewAuditHandler(deps.AuditService)

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/readyz", healthHandler.Readyz)

	if deps.Config.HTTP.EnableMetrics {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
	if deps.Config.HTTP.EnableSwagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := r.Group(deps.Config.HTTP.BasePath)
	{
		auth := api.Group("/auth")
		{
			auth.POST("/request-otp", authHandler.RequestOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password/request-otp", authHandler.ForgotPasswordRequestOTP)
			auth.POST("/forgot-password/confirm", authHandler.ForgotPasswordConfirm)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		me := api.Group("/me")
		me.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			me.GET("", userHandler.Me)
			me.PATCH("/profile", userHandler.UpdateProfile)
			me.PATCH("/preferences", userHandler.UpdatePreferences)
			me.POST("/device-tokens", userHandler.SaveDeviceToken)
		}

		caregivers := api.Group("/caregivers")
		caregivers.Use(middleware.RequireAuth(deps.Config.JWT))
		caregivers.Use(middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin))
		{
			caregivers.POST("/assignments", caregiverHandler.CreateAssignment)
			caregivers.GET("/assignments", caregiverHandler.ListAssignments)
		}

		medicines := api.Group("/medicines")
		medicines.Use(middleware.RequireAuth(deps.Config.JWT))
		medicines.GET("/categories", medicineHandler.ListCategories)
		medicines.GET("/categories/:id/items", medicineHandler.ListCategoryItems)
		medicines.GET("/dosage-options", medicineHandler.GetDosageOptions)
		medicines.GET("/meal-timing-options", medicineHandler.GetMealTimingOptions)
		medicines.GET("/master", medicineHandler.ListMaster)
		medicines.Use(middleware.RequireRoles(constants.RolePatient, constants.RoleNurse, constants.RoleAdmin))
		{
			medicines.POST("/patient", medicineHandler.CreatePatientMedicine)
			medicines.GET("/patient", medicineHandler.ListPatientMedicines)
			medicines.PATCH("/patient/:id", medicineHandler.UpdatePatientMedicine)
			medicines.DELETE("/patient/:id", medicineHandler.DeletePatientMedicine)
			medicines.POST("/patient/:id/schedules", medicineHandler.CreateSchedule)
			medicines.DELETE("/schedules/:id", medicineHandler.DeleteSchedule)
		}

		intake := api.Group("/intake")
		intake.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			intake.POST("", middleware.RequireRoles(constants.RolePatient, constants.RoleNurse, constants.RoleAdmin), intakeHandler.CreateIntake)
			intake.GET("/history", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), intakeHandler.ListHistory)
		}

		health := api.Group("/health")
		health.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			health.POST("/records", middleware.RequireRoles(constants.RolePatient, constants.RoleNurse, constants.RoleAdmin), healthRecordHandler.CreateHealthRecord)
			health.GET("/records", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), healthRecordHandler.ListHealthRecords)
		}

		assessments := api.Group("/assessments")
		assessments.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			assessments.POST("/daily", middleware.RequireRoles(constants.RolePatient, constants.RoleNurse, constants.RoleAdmin), healthRecordHandler.CreateDailyAssessment)
			assessments.GET("/daily", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), healthRecordHandler.ListDailyAssessments)
		}

		appointments := api.Group("/appointments")
		appointments.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			appointments.GET("", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), appointmentHandler.ListAppointments)
			appointments.POST("", middleware.RequireRoles(constants.RolePatient, constants.RoleNurse, constants.RoleAdmin), appointmentHandler.CreateAppointment)
			appointments.PATCH("/:id/status", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), appointmentHandler.UpdateStatus)
			appointments.DELETE("/:id", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), appointmentHandler.DeleteAppointment)
			appointments.POST("/:id/notes", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), appointmentHandler.CreateNurseVisitNote)
		}

		visits := api.Group("/visits")
		visits.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			visits.GET("/history", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), appointmentHandler.ListVisitHistory)
		}

		content := api.Group("/content")
		content.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			content.GET("/health/categories", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), contentHandler.ListHealthCategories)
			content.GET("/health", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), contentHandler.ListHealthContent)
			content.POST("/health", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), contentHandler.CreateHealthContent)
			content.PATCH("/health/:id", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), contentHandler.UpdateHealthContent)
			content.POST("/health/:id/publish", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), contentHandler.PublishHealthContent)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.RequireAuth(deps.Config.JWT))
		{
			notifications.GET("/upcoming", middleware.RequireRoles(constants.RolePatient, constants.RoleCaregiver, constants.RoleNurse, constants.RoleAdmin), notificationHandler.ListUpcoming)
		}

		support := api.Group("/support")
		{
			support.GET("/emergency", supportHandler.EmergencyInfo)
			chat := support.Group("/chat")
			chat.Use(middleware.RequireAuth(deps.Config.JWT))
			chat.POST("/requests", middleware.RequireRoles(constants.RolePatient), supportHandler.CreateChatRequest)
			chat.GET("/requests", middleware.RequireRoles(constants.RoleNurse, constants.RoleAdmin), supportHandler.ListChatRequests)
		}

		staff := api.Group("/staff")
		{
			staff.POST("/login", adminHandler.StaffLogin)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(deps.Config.JWT))
		admin.Use(middleware.RequireRoles(constants.RoleAdmin))
		{
			admin.GET("/patients", adminHandler.ListPatients)
			admin.GET("/patients/:id", adminHandler.GetPatient)
			admin.GET("/adherence", adminHandler.ListAdherence)
			admin.GET("/audit-logs", auditHandler.ListAuditLogs)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		httpx.Fail(c, domain.NewError(constants.UserNotFound, "not found"))
	})

	return r
}
