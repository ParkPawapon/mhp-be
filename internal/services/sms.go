package services

import (
	"context"

	"go.uber.org/zap"
)

type SmsSender interface {
	SendOTP(ctx context.Context, phone, otpCode, refCode string) error
}

type ConsoleSender struct {
	Logger *zap.Logger
}

func (s ConsoleSender) SendOTP(ctx context.Context, phone, otpCode, refCode string) error {
	_ = ctx
	if s.Logger != nil {
		s.Logger.Info("otp", zap.String("phone", phone), zap.String("otp_code", otpCode), zap.String("ref_code", refCode))
	}
	return nil
}
