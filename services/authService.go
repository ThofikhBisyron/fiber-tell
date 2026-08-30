package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"tell-be/lib"
	"tell-be/repositories"
	"time"
)

type AuthService struct {
	userRepo     *repositories.UserRepo
	emailOtpRepo *repositories.OTPRepo
	emailService *EmailService
}

func NewAuthService(
	userRepo *repositories.UserRepo,
	emailOtpRepo *repositories.OTPRepo,
	emailService *EmailService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailOtpRepo: emailOtpRepo,
		emailService: emailService,
	}
}

func generateOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *AuthService) RequestOtp(
	ctx context.Context,
	email string,
) error {
	email = strings.ToLower(strings.TrimSpace(email))

	code := generateOtp()

	codeHash, err := lib.HashPassword(code)

	if err != nil {
		return err
	}

	expires_at := time.Now().Add(5 * time.Minute)

	if err := s.emailOtpRepo.DeleteOtpByEmail(
		ctx,
		email,
	); err != nil {
		return err
	}

	if err := s.emailOtpRepo.CreateOtp(
		ctx,
		email,
		codeHash,
		expires_at,
	); err != nil {
		return err
	}

	if err := s.emailService.SendOtp(
		email,
		code,
	); err != nil {
		return err
	}

	return nil

}

func (s *AuthService) VerifyOtp(
	ctx context.Context,
	email string,
	code string,
) error {
	email = strings.ToLower(strings.TrimSpace(email))

	otp, err := s.emailOtpRepo.FindOtpByEmail(
		ctx,
		email,
	)

	if err != nil {
		return fmt.Errorf("Otp Not Found")
	}

	if time.Now().After(otp.Expires_at) {
		_ = s.emailOtpRepo.DeleteOtpById(ctx, otp.Id)

		return fmt.Errorf("otp expired, please send code again")
	}

	if otp.Attempts >= 5 {
		_ = s.emailOtpRepo.DeleteOtpById(ctx, otp.Id)

		return fmt.Errorf("too many attempts")
	}

	if !lib.VerifyPassword(code, otp.Code_hash) {
		if err := s.emailOtpRepo.IncreaseAttempts(
			ctx,
			int(otp.Id),
		); err != nil {
			return err
		}

		return fmt.Errorf("Invalid Otp")
	}

	if err := s.emailOtpRepo.DeleteOtpById(
		ctx,
		otp.Id,
	); err != nil {
		return err
	}

	user, err := s.userRepo.FindUserByEmail(
		ctx,
		email,
	)

	if err != nil {
		user, err = s.userRepo.CreateUser(
			ctx,
			email,
		)

		if err != nil {
			return err
		}
	}

	fmt.Println("Login user:", user.Id)

	// TODO:
	// Generate JWT access token
	// Generate JWT refresh token

	return nil
}
