package constants

const (
	MealTimingBeforeMeal           = "BEFORE_MEAL"
	MealTimingAfterMeal            = "AFTER_MEAL"
	MealTimingAfterMealImmediately = "AFTER_MEAL_IMMEDIATELY"
	MealTimingBeforeBed            = "BEFORE_BED"
	MealTimingUntilFinished        = "UNTIL_FINISHED"
	MealTimingNoMilk               = "NO_MILK"
	MealTimingOther                = "OTHER"
)

var MealTimingOptions = []string{
	MealTimingBeforeMeal,
	MealTimingAfterMeal,
	MealTimingAfterMealImmediately,
	MealTimingBeforeBed,
	MealTimingUntilFinished,
	MealTimingNoMilk,
	MealTimingOther,
}

var DosageOptions = []string{"1/4", "1/2", "1", "2"}

const (
	SupportCategoryGeneral     = "GENERAL"
	SupportCategoryMedicine    = "MEDICINE"
	SupportCategoryAppointment = "APPOINTMENT"
	SupportCategoryTech        = "TECH"
)

var SupportCategories = []string{
	SupportCategoryGeneral,
	SupportCategoryMedicine,
	SupportCategoryAppointment,
	SupportCategoryTech,
}

const (
	ContentCategoryHypertensionKnowledge = "HYPERTENSION_KNOWLEDGE"
	ContentCategoryHypertensionControl   = "HYPERTENSION_CONTROL"
	ContentCategoryAbnormalSymptoms      = "ABNORMAL_SYMPTOMS"
	ContentCategoryCvRiskScore           = "CV_RISK_SCORE"
	ContentCategoryFoodAndMedicine       = "FOOD_AND_MEDICINE"
	ContentCategoryExerciseAndBMI        = "EXERCISE_AND_BMI"
	ContentCategoryStressManagement      = "STRESS_MANAGEMENT"
	ContentCategorySleep                 = "SLEEP"
)

var HealthContentCategories = []string{
	ContentCategoryHypertensionKnowledge,
	ContentCategoryHypertensionControl,
	ContentCategoryAbnormalSymptoms,
	ContentCategoryCvRiskScore,
	ContentCategoryFoodAndMedicine,
	ContentCategoryExerciseAndBMI,
	ContentCategoryStressManagement,
	ContentCategorySleep,
}
