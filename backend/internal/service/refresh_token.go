package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
	"github.com/sqlc-dev/pqtype"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)

type RefreshTokenService struct {
	db   *sql.DB
	repo repository.Querier
}

func NewRefreshTokenService(db *sql.DB, repo repository.Querier) *RefreshTokenService {
	return &RefreshTokenService{db: db, repo: repo}
}

type RefreshToken struct {
	ID            uuid.UUID
	GrantWriterID uuid.UUID
	Token         string
	UserAgent     string
	IpAddress     string
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

type CreateTokenInput struct {
	GrantWriterID uuid.UUID
	UserAgent     string
	IpAddress     string
}

type RotateTokenInput struct {
	UserAgent string
	IpAddress string
}

func (s *RefreshTokenService) CreateToken(ctx context.Context, input CreateTokenInput) (*RefreshToken, error) {
	newTokenValue, newTokenHash := generateRefreshToken()

	record, err := s.repo.CreateToken(ctx, repository.CreateTokenParams{
		GrantWriterID: input.GrantWriterID,
		Token:         newTokenHash,
		IpAddress:     toInetFromString(input.IpAddress),
		UserAgent:     toNullString(input.UserAgent),
	})
	if err != nil {
		return nil, err
	}

	newTokenObject := toRefreshTokenResponse(record)

	newTokenObject.Token = newTokenValue

	return newTokenObject, nil
}

func (s *RefreshTokenService) RotateToken(ctx context.Context, tokenValue string, input RotateTokenInput) (*RefreshToken, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txRepo := repository.New(tx)
	tokenHash := hashToken(tokenValue)

	oldRecord, err := txRepo.GetRefreshTokenByTokenValue(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	if time.Now().After(oldRecord.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	err = txRepo.DeleteRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	newTokenValue, newTokenHash := generateRefreshToken()
	record, err := txRepo.CreateToken(ctx, repository.CreateTokenParams{
		GrantWriterID: oldRecord.GrantWriterID,
		Token:         newTokenHash,
		IpAddress:     toInetFromString(input.IpAddress),
		UserAgent:     toNullString(input.UserAgent),
	})
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	newTokenObject := toRefreshTokenResponse(record)
	newTokenObject.Token = newTokenValue

	return newTokenObject, nil
}

func (s *RefreshTokenService) CountValidTokens(ctx context.Context, grantWriterID uuid.UUID) (int64, error) {
	tokenCount, err := s.repo.CountValidTokens(ctx, grantWriterID)
	if err != nil {
		return -1, err
	}

	return tokenCount, nil
}

func (s *RefreshTokenService) DeleteRefreshToken(ctx context.Context, tokenValue string) error {
	tokenHash := hashToken(tokenValue)
	err := s.repo.DeleteRefreshToken(ctx, tokenHash)

	return err
}

func (s *RefreshTokenService) DeleteAllRefreshTokens(ctx context.Context, grantWriterID uuid.UUID) error {
	err := s.repo.DeleteAllRefreshTokens(ctx, grantWriterID)

	return err
}

func (s *RefreshTokenService) DeleteExpiredTokens(ctx context.Context) error {
	err := s.repo.DeleteExpiredTokens(ctx)

	return err
}

func generateRefreshToken() (string, string) {
	newToken := rand.Text()
	hash := sha256.Sum256([]byte(newToken))
	return newToken, hex.EncodeToString(hash[:])
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func toRefreshTokenResponse(r repository.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:            r.ID,
		GrantWriterID: r.GrantWriterID,
		UserAgent:     r.UserAgent.String,
		IpAddress:     r.IpAddress.IPNet.IP.String(),
		CreatedAt:     r.CreatedAt,
		ExpiresAt:     r.ExpiresAt,
	}
}

func toInetFromString(ipAddressString string) pqtype.Inet {
	host, _, err := net.SplitHostPort(ipAddressString)

	if err != nil {
		host = ipAddressString
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return pqtype.Inet{Valid: false}
	}

	var mask net.IPMask

	if ip.To4() != nil {
		mask = ip.DefaultMask()
	} else {
		mask = net.CIDRMask(128, 128)
	}

	return pqtype.Inet{
		IPNet: net.IPNet{
			IP:   ip,
			Mask: mask,
		},
		Valid: true,
	}
}
