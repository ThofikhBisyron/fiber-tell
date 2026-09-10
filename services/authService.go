package services

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"tell-be/lib"
	"tell-be/models"
	"tell-be/repositories"
	"time"
)

type LoginResult struct {
	User         models.User
	AccessToken  string
	RefreshToken string
}
type AuthService struct {
	userRepo     *repositories.UserRepo
	profileRepo  *repositories.ProfileRepo
	emailOtpRepo *repositories.OTPRepo
	emailService *EmailService
	jwtService   *JwtService
}

func NewAuthService(
	userRepo *repositories.UserRepo,
	profileRepo *repositories.ProfileRepo,
	emailOtpRepo *repositories.OTPRepo,
	emailService *EmailService,
	jwtService *JwtService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		emailOtpRepo: emailOtpRepo,
		emailService: emailService,
		jwtService:   jwtService,
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
) (*LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	otp, err := s.emailOtpRepo.FindOtpByEmail(
		ctx,
		email,
	)

	if err != nil {
		return nil, fmt.Errorf("Otp Not Found")
	}

	if time.Now().After(otp.Expires_at) {
		_ = s.emailOtpRepo.DeleteOtpById(ctx, otp.Id)

		return nil, fmt.Errorf("otp expired, please send code again")
	}

	if otp.Attempts >= 5 {
		_ = s.emailOtpRepo.DeleteOtpById(ctx, otp.Id)

		return nil, fmt.Errorf("too many attempts")
	}

	if !lib.VerifyPassword(code, otp.Code_hash) {
		if err := s.emailOtpRepo.IncreaseAttempts(
			ctx,
			int(otp.Id),
		); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Invalid Otp")
	}

	if err := s.emailOtpRepo.DeleteOtpById(
		ctx,
		otp.Id,
	); err != nil {
		return nil, err
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
			return nil, err
		}

		_, err = s.profileRepo.CreateProfile(
			ctx,
			user.Id,
		)

		if err != nil {
			return nil, err
		}
	}

	fmt.Println("Login user:", user.Id)

	accessToken, err := s.jwtService.GenerateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
	// TODO:
	// Generate JWT access token
	// Generate JWT refresh token

}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*LoginResult, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	userIDString, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := s.userRepo.FindUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.Id)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
