package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/cache"
	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/database/postgres"
	"github.com/ParkPawapon/mhp-be/internal/logging"
	"github.com/ParkPawapon/mhp-be/internal/observability"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
	"github.com/ParkPawapon/mhp-be/internal/server"
	"github.com/ParkPawapon/mhp-be/internal/services"
	httptransport "github.com/ParkPawapon/mhp-be/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	logger, err := logging.New(cfg.App.LogLevel)
	if err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
	defer logger.Sync()

	shutdownTracer, err := observability.InitTracerProvider(cfg.Observability)
	if err != nil {
		logger.Fatal("otel init failed", zap.Error(err))
	}
	defer func() {
		_ = shutdownTracer(context.Background())
	}()

	observability.InitMetrics()

	db, err := postgres.New(cfg.DB)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}

	redisClient, err := cache.New(cfg.Redis)
	if err != nil {
		logger.Fatal("redis connection failed", zap.Error(err))
	}

	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	profileRepo := repositories.NewProfileRepository(db)
	caregiverRepo := repositories.NewCaregiverRepository(db)

	smsSender, err := newSmsSender(cfg, logger)
	if err != nil {
		logger.Fatal("sms sender init failed", zap.Error(err))
	}
	authService := services.NewAuthService(cfg, authRepo, userRepo, redisClient, smsSender)
	userService := services.NewUserService(userRepo, profileRepo, redisClient)
	caregiverService := services.NewCaregiverService(caregiverRepo)

	router := httptransport.NewRouter(httptransport.Dependencies{
		Config:             cfg,
		Logger:             logger,
		DB:                 db,
		Redis:              redisClient,
		AuthService:        authService,
		UserService:        userService,
		CaregiverService:   caregiverService,
		MedicineService:    services.NewMedicineService(),
		IntakeService:      services.NewIntakeService(),
		HealthService:      services.NewHealthService(),
		AppointmentService: services.NewAppointmentService(),
		ContentService:     services.NewContentService(),
		AdminService:       services.NewAdminService(),
		AuditService:       services.NewAuditService(),
	})

	addr := server.Address(cfg.HTTP.Host, cfg.HTTP.Port)
	srv := server.New(addr, router, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout, cfg.HTTP.IdleTimeout)

	go func() {
		logger.Info("server started", zap.String("addr", addr))
		var err error
		if cfg.HTTP.TLSCertFile != "" && cfg.HTTP.TLSKeyFile != "" {
			err = srv.StartTLS(cfg.HTTP.TLSCertFile, cfg.HTTP.TLSKeyFile)
		} else {
			err = srv.Start()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}
	logger.Info("server stopped")
}

func newSmsSender(cfg config.Config, logger *zap.Logger) (services.SmsSender, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.SMS.Provider))
	switch provider {
	case "", "console":
		return services.ConsoleSender{Logger: logger}, nil
	case "disabled", "none":
		return nil, nil
	case "thaibulksms":
		return services.NewThaiBulkSMSSender(cfg.SMS.ThaiBulkSMS, logger)
	default:
		return nil, fmt.Errorf("unsupported sms provider: %s", cfg.SMS.Provider)
	}
}
