package constants

type NotificationStatus string

const (
	NotificationPending   NotificationStatus = "PENDING"
	NotificationSent      NotificationStatus = "SENT"
	NotificationCancelled NotificationStatus = "CANCELLED"
	NotificationFailed    NotificationStatus = "FAILED"
)

const (
	TemplateMedBeforeMeal5Min  = "MED_BEFORE_MEAL_5MIN"
	TemplateMedBeforeMeal20Min = "MED_BEFORE_MEAL_20MIN"
	TemplateMedAfterMealNow    = "MED_AFTER_MEAL_NOW"
	TemplateAppt5Days          = "APPT_5D"
	TemplateAppt1Day           = "APPT_1D"
	TemplateWeeklyHealthLog    = "WEEKLY_HEALTH_LOG"
)
