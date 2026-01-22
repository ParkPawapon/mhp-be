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
	if err := seedMedicineCategories(ctx, dbConn); err != nil {
		log.Fatalf("seed medicine categories failed: %v", err)
	}
	if err := seedNotificationTemplates(ctx, dbConn); err != nil {
		log.Fatalf("seed notification templates failed: %v", err)
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

func seedMedicineCategories(ctx context.Context, dbConn *gorm.DB) error {
	category := db.MedicineCategory{
		Name:      "Hypertension",
		Code:      "HYPERTENSION",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}

	var existing db.MedicineCategory
	err := dbConn.WithContext(ctx).Where("code = ?", category.Code).First(&existing).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := dbConn.WithContext(ctx).Create(&category).Error; err != nil {
			return err
		}
		existing = category
	}

	items := []db.MedicineCategoryItem{
		{CategoryID: existing.ID, DisplayName: "Beta-blocker: Atenolol", DefaultDosageText: strPtr("1"), IsActive: true, CreatedAt: time.Now().UTC()},
		{CategoryID: existing.ID, DisplayName: "Beta-blocker: Propranolol", DefaultDosageText: strPtr("1"), IsActive: true, CreatedAt: time.Now().UTC()},
		{CategoryID: existing.ID, DisplayName: "Calcium channel blocker: Amlodipine 5 mg", DefaultDosageText: strPtr("1"), IsActive: true, CreatedAt: time.Now().UTC()},
		{CategoryID: existing.ID, DisplayName: "Calcium channel blocker: Amlodipine 10 mg", DefaultDosageText: strPtr("1"), IsActive: true, CreatedAt: time.Now().UTC()},
	}

	for _, item := range items {
		var count int64
		if err := dbConn.WithContext(ctx).Model(&db.MedicineCategoryItem{}).
			Where("category_id = ? AND display_name = ?", item.CategoryID, item.DisplayName).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := dbConn.WithContext(ctx).Create(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedNotificationTemplates(ctx context.Context, dbConn *gorm.DB) error {
	templates := []db.NotificationTemplate{
		{Code: constants.TemplateMedBeforeMeal5Min, Title: "Medicine reminder", Body: "Reminder: take your medicine before meal in 5 minutes.", IsActive: true, CreatedAt: time.Now().UTC()},
		{Code: constants.TemplateMedBeforeMeal20Min, Title: "Medicine reminder", Body: "Reminder: please take your medicine if you have not yet.", IsActive: true, CreatedAt: time.Now().UTC()},
		{Code: constants.TemplateMedAfterMealNow, Title: "Medicine reminder", Body: "หลังทานอาหารอย่าลืมทานยา", IsActive: true, CreatedAt: time.Now().UTC()},
		{Code: constants.TemplateAppt5Days, Title: "Appointment reminder", Body: "Upcoming appointment in 5 days.", IsActive: true, CreatedAt: time.Now().UTC()},
		{Code: constants.TemplateAppt1Day, Title: "Appointment reminder", Body: "Upcoming appointment tomorrow.", IsActive: true, CreatedAt: time.Now().UTC()},
		{Code: constants.TemplateWeeklyHealthLog, Title: "Weekly health log", Body: "Please complete your weekly health behavior log.", IsActive: true, CreatedAt: time.Now().UTC()},
	}

	for _, tpl := range templates {
		var count int64
		if err := dbConn.WithContext(ctx).Model(&db.NotificationTemplate{}).Where("code = ?", tpl.Code).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := dbConn.WithContext(ctx).Create(&tpl).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func strPtr(v string) *string {
	return &v
}
