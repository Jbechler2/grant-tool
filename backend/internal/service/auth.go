package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already in use")
	ErrTooManySessions    = errors.New("max sessions exceeded")
)

type AuthService struct {
	repo                repository.Querier
	jwtSecret           string
	jwtExpiry           time.Duration
	refreshTokenService *RefreshTokenService
}

func NewAuthService(repo repository.Querier, jwtSecret string, jwtExpiryMinutes int, refreshTokenService *RefreshTokenService) *AuthService {
	return &AuthService{
		repo:                repo,
		jwtSecret:           jwtSecret,
		jwtExpiry:           time.Duration(jwtExpiryMinutes) * time.Minute,
		refreshTokenService: refreshTokenService,
	}
}

type RegisterInput struct {
	Email     string
	Password  string
	UserAgent string
	IpAddress string
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IpAddress string
}

type AuthResult struct {
	Token         string
	RefreshToken  string
	RefreshExpiry time.Time
	User          repository.User
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	existing, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err == nil && existing.ID != uuid.Nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         repository.UserRoleGrantWriter,
	})
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.refreshTokenService.CreateToken(ctx, CreateTokenInput{
		GrantWriterID: user.ID,
		UserAgent:     input.UserAgent,
		IpAddress:     input.IpAddress,
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, RefreshToken: newRefreshToken.Token, RefreshExpiry: newRefreshToken.ExpiresAt, User: user}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		_ = err
	}

	validTokens, err := s.refreshTokenService.CountValidTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if validTokens >= 8 {
		return nil, ErrTooManySessions
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.refreshTokenService.CreateToken(ctx, CreateTokenInput{
		GrantWriterID: user.ID,
		UserAgent:     input.UserAgent,
		IpAddress:     input.IpAddress,
	})
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, RefreshToken: newRefreshToken.Token, RefreshExpiry: newRefreshToken.ExpiresAt, User: user}, nil
}

func (s *AuthService) generateToken(user repository.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.jwtExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) RotateToken(ctx context.Context, tokenValue string, input RotateTokenInput) (*AuthResult, error) {
	newToken, err := s.refreshTokenService.RotateToken(ctx, tokenValue, input)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, newToken.GrantWriterID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token:         accessToken,
		RefreshToken:  newToken.Token,
		RefreshExpiry: newToken.ExpiresAt,
		User:          user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, tokenValue string) error {
	return s.refreshTokenService.DeleteRefreshToken(ctx, tokenValue)
}
