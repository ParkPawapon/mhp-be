CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$
BEGIN
    CREATE TYPE role_type AS ENUM ('PATIENT', 'NURSE', 'ADMIN', 'CAREGIVER');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE gender_type AS ENUM ('MALE', 'FEMALE', 'OTHER');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE med_intake_status AS ENUM ('TAKEN', 'MISSED', 'SKIPPED');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE appt_category AS ENUM ('HOSPITAL', 'HOME_VISIT');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE appt_status_type AS ENUM ('PENDING', 'CONFIRMED', 'COMPLETED', 'CANCELLED');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role role_type NOT NULL DEFAULT 'PATIENT',
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    line_user_id VARCHAR(100) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    hn VARCHAR(20) UNIQUE,
    citizen_id VARCHAR(13) UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender gender_type,
    blood_type VARCHAR(5),
    address_text TEXT,
    gps_lat DECIMAL(10,8),
    gps_long DECIMAL(11,8),
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(15),
    avatar_url TEXT
);

CREATE TABLE IF NOT EXISTS caregiver_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES users(id),
    caregiver_id UUID NOT NULL REFERENCES users(id),
    relationship VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth_otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    otp_code VARCHAR(10) NOT NULL,
    ref_code VARCHAR(10) NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    is_used BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medicines_master (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_name VARCHAR(255) NOT NULL,
    generic_name VARCHAR(255),
    dosage_unit VARCHAR(50) NOT NULL,
    default_image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS patient_medicines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    medicine_master_id UUID REFERENCES medicines_master(id),
    custom_name VARCHAR(255),
    dosage_amount VARCHAR(100) NOT NULL,
    instruction TEXT,
    indication TEXT,
    my_drug_image_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS medicine_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_medicine_id UUID NOT NULL REFERENCES patient_medicines(id) ON DELETE CASCADE,
    time_slot TIME NOT NULL,
    meal_timing VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS intake_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    schedule_id UUID REFERENCES medicine_schedules(id),
    target_date DATE NOT NULL,
    taken_at TIMESTAMPTZ,
    status med_intake_status NOT NULL,
    skip_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS health_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    record_date DATE NOT NULL,
    time_period VARCHAR(20) NOT NULL,
    systolic_bp INTEGER,
    diastolic_bp INTEGER,
    pulse_rate INTEGER,
    weight_kg DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS daily_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_date DATE NOT NULL,
    exercise_minutes INTEGER DEFAULT 0,
    sleep_quality VARCHAR(50),
    stress_level INTEGER,
    diet_compliance VARCHAR(50),
    symptoms JSONB,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    creator_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    appt_type appt_category NOT NULL,
    appt_datetime TIMESTAMPTZ NOT NULL,
    location_name VARCHAR(255),
    slip_image_url TEXT,
    status appt_status_type DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS nurse_visit_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    nurse_id UUID NOT NULL REFERENCES users(id),
    visit_details TEXT NOT NULL,
    vital_signs_summary JSONB,
    next_action_plan TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS health_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    body_content TEXT,
    thumbnail_url TEXT,
    external_video_url TEXT,
    category VARCHAR(50),
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id),
    target_user_id UUID REFERENCES users(id),
    action_type VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_caregiver_assignments_patient_id ON caregiver_assignments(patient_id);
CREATE INDEX IF NOT EXISTS idx_caregiver_assignments_caregiver_id ON caregiver_assignments(caregiver_id);
CREATE INDEX IF NOT EXISTS idx_auth_otp_codes_phone_number ON auth_otp_codes(phone_number);
CREATE INDEX IF NOT EXISTS idx_patient_medicines_user_id ON patient_medicines(user_id);
CREATE INDEX IF NOT EXISTS idx_medicine_schedules_patient_medicine_id ON medicine_schedules(patient_medicine_id);
CREATE INDEX IF NOT EXISTS idx_intake_history_user_id ON intake_history(user_id);
CREATE INDEX IF NOT EXISTS idx_intake_history_schedule_id ON intake_history(schedule_id);
CREATE INDEX IF NOT EXISTS idx_health_records_user_id ON health_records(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_assessments_user_id ON daily_assessments(user_id);
CREATE INDEX IF NOT EXISTS idx_appointments_user_id ON appointments(user_id);
CREATE INDEX IF NOT EXISTS idx_nurse_visit_notes_appointment_id ON nurse_visit_notes(appointment_id);
CREATE INDEX IF NOT EXISTS idx_nurse_visit_notes_nurse_id ON nurse_visit_notes(nurse_id);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_profiles_hn ON user_profiles(hn);
CREATE INDEX IF NOT EXISTS idx_user_profiles_citizen_id ON user_profiles(citizen_id);
CREATE INDEX IF NOT EXISTS idx_intake_history_user_id_target_date ON intake_history(user_id, target_date);
CREATE INDEX IF NOT EXISTS idx_appointments_user_id_appt_datetime ON appointments(user_id, appt_datetime);
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp_actor_id ON audit_logs(timestamp, actor_id);
