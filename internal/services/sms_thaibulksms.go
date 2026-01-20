package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/config"
)

type ThaiBulkSMSSender struct {
	cfg    config.ThaiBulkSMSConfig
	client *http.Client
	logger *zap.Logger
}

func NewThaiBulkSMSSender(cfg config.ThaiBulkSMSConfig, logger *zap.Logger) (*ThaiBulkSMSSender, error) {
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.Endpoint = strings.TrimSpace(cfg.Endpoint)
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.APISecret = strings.TrimSpace(cfg.APISecret)
	cfg.AuthMode = strings.ToLower(strings.TrimSpace(cfg.AuthMode))

	if cfg.BaseURL == "" {
		return nil, errors.New("thaibulksms base url required")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("thaibulksms endpoint required")
	}
	if !strings.HasPrefix(cfg.Endpoint, "/") {
		cfg.Endpoint = "/" + cfg.Endpoint
	}
	if cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, errors.New("thaibulksms api key/secret required")
	}
	if cfg.AuthMode == "" {
		cfg.AuthMode = "basic"
	}
	if cfg.AuthMode != "basic" && cfg.AuthMode != "body" {
		return nil, errors.New("thaibulksms auth mode must be basic or body")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &ThaiBulkSMSSender{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
		logger: logger,
	}, nil
}

func (s *ThaiBulkSMSSender) SendOTP(ctx context.Context, phone, otpCode, refCode string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return errors.New("phone required")
	}

	message := s.renderTemplate(otpCode, refCode)
	form := url.Values{}
	form.Set("msisdn", phone)
	form.Set("message", message)
	if strings.TrimSpace(s.cfg.SenderID) != "" {
		form.Set("sender", strings.TrimSpace(s.cfg.SenderID))
	}
	if s.cfg.AuthMode == "body" {
		form.Set("key", s.cfg.APIKey)
		form.Set("secret", s.cfg.APISecret)
	}

	urlStr := strings.TrimRight(s.cfg.BaseURL, "/") + s.cfg.Endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if s.cfg.AuthMode == "basic" {
		req.SetBasicAuth(s.cfg.APIKey, s.cfg.APISecret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logError("thaibulksms request failed", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		err := fmt.Errorf("thaibulksms error status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
		s.logError("thaibulksms response failed", err)
		return err
	}

	return nil
}

func (s *ThaiBulkSMSSender) renderTemplate(otpCode, refCode string) string {
	tpl := strings.TrimSpace(s.cfg.OTPTemplate)
	if tpl == "" {
		tpl = "Your OTP is {{otp}} (ref: {{ref}})"
	}
	msg := strings.ReplaceAll(tpl, "{{otp}}", otpCode)
	msg = strings.ReplaceAll(msg, "{{ref}}", refCode)
	return msg
}

func (s *ThaiBulkSMSSender) logError(msg string, err error) {
	if s.logger == nil {
		return
	}
	s.logger.Error(msg, zap.Error(err), zap.String("provider", "thaibulksms"))
}
