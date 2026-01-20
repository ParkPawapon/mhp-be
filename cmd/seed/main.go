package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/database/postgres"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	dbConn, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	ctx := context.Background()
	if err := seedUsers(ctx, dbConn); err != nil {
		log.Fatalf("seed users failed: %v", err)
	}
	if err := seedMedicines(ctx, dbConn); err != nil {
		log.Fatalf("seed medicines failed: %v", err)
	}

	log.Println("seed completed")
}

func seedUsers(ctx context.Context, dbConn *gorm.DB) error {
	adminPwd, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	nursePwd, _ := bcrypt.GenerateFromPassword([]byte("Nurse123!"), bcrypt.DefaultCost)

	users := []db.User{
		{Username: "admin", PasswordHash: string(adminPwd), Role: constants.RoleAdmin, IsActive: true, IsVerified: true},
		{Username: "nurse", PasswordHash: string(nursePwd), Role: constants.RoleNurse, IsActive: true, IsVerified: true},
	}

	for _, user := range users {
		var count int64
		if err := dbConn.WithContext(ctx).Model(&db.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := dbConn.WithContext(ctx).Create(&user).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedMedicines(ctx context.Context, dbConn *gorm.DB) error {
	meds := []db.MedicineMaster{
		{TradeName: "Paracetamol", GenericName: strPtr("Acetaminophen"), DosageUnit: "mg", CreatedAt: time.Now().UTC()},
		{TradeName: "Amoxicillin", GenericName: strPtr("Amoxicillin"), DosageUnit: "mg", CreatedAt: time.Now().UTC()},
	}

	for _, med := range meds {
		var count int64
		if err := dbConn.WithContext(ctx).Model(&db.MedicineMaster{}).Where("trade_name = ?", med.TradeName).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := dbConn.WithContext(ctx).Create(&med).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func strPtr(v string) *string {
	return &v
}
