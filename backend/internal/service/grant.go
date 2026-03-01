package service

import (
	"context"
	"errors"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
)

var (
	ErrGrantNotFound     = errors.New("grant not found")
	ErrGrantUnauthorized = errors.New("unauthorized access to grant")
)

type GrantService struct {
	repo repository.Querier
}

func NewGrantService(repo repository.Querier) *GrantService {
	return &GrantService{repo: repo}
}

type Grant struct {
	ID                        uuid.UUID
	GrantWriterID             uuid.UUID
	Title                     string
	FunderName                string
	FunderWebsite             string
	Description               string
	AwardAmountMin            *float64
	AwardAmountMax            *float64
	EligibilityNotes          string
	EstimatedApplicationHours *float64
	Visibility                string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type CreateGrantInput struct {
	GrantWriterID             uuid.UUID
	Title                     string
	FunderName                string
	FunderWebsite             *string
	Description               *string
	AwardAmountMin            *float64
	AwardAmountMax            *float64
	EligibilityNotes          *string
	EstimatedApplicationHours *float64
	Visibility                *string
}

type Deadline struct {
	ID          uuid.UUID `json:"id"`
	GrantID     uuid.UUID `json:"grant_id"`
	Label       string    `json:"label"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type AddDeadlineInput struct {
	GrantWriterID uuid.UUID
	GrantID       uuid.UUID
	Label         string
	Date          time.Time
	Description   *string
}

type UpdateGrantInput struct {
	ID                        uuid.UUID
	GrantWriterID             uuid.UUID
	Title                     *string
	FunderName                *string
	FunderWebsite             *string
	Description               *string
	AwardAmountMin            *float64
	AwardAmountMax            *float64
	EligibilityNotes          *string
	EstimatedApplicationHours *float64
	Visibility                *string
}

func (s *GrantService) CreateGrant(ctx context.Context, input CreateGrantInput) (*Grant, error) {
	record, err := s.repo.CreateGrant(ctx, repository.CreateGrantParams{
		Title:                     input.Title,
		GrantWriterID:             input.GrantWriterID,
		FunderName:                input.FunderName,
		Visibility:                repository.GrantVisibilityPrivate,
		FunderWebsite:             toNullStringFromPtr(input.FunderWebsite),
		Description:               toNullStringFromPtr(input.Description),
		AwardAmountMin:            float64ToNullString(input.AwardAmountMin),
		AwardAmountMax:            float64ToNullString(input.AwardAmountMax),
		EligibilityNotes:          toNullStringFromPtr(input.EligibilityNotes),
		EstimatedApplicationHours: float64ToNullString(input.EstimatedApplicationHours),
	})
	if err != nil {
		return nil, err
	}

	return toGrantResponse(record), nil
}

func (s *GrantService) GetGrantByID(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID) (*Grant, error) {
	record, err := s.repo.GetGrantByID(ctx, repository.GetGrantByIDParams{
		GrantWriterID: grantWriterID,
		ID:            grantID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}

	return toGrantResponse(record), nil
}

func (s *GrantService) GetAllGrants(ctx context.Context, grantWriterID uuid.UUID) ([]Grant, error) {
	records, err := s.repo.GetAllGrants(ctx, grantWriterID)
	if err != nil {
		return nil, err
	}

	grants := make([]Grant, len(records))
	for i, record := range records {
		grants[i] = *toGrantResponse(record)
	}

	return grants, nil
}

func (s *GrantService) UpdateGrant(ctx context.Context, input UpdateGrantInput) (*Grant, error) {
	existingRecord, err := s.repo.GetGrantByID(ctx, repository.GetGrantByIDParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.ID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}

	record, err := s.repo.UpdateGrant(ctx, repository.UpdateGrantParams{
		GrantWriterID:             input.GrantWriterID,
		ID:                        input.ID,
		Title:                     mergeString(input.Title, existingRecord.Title),
		FunderName:                mergeString(input.FunderName, existingRecord.FunderName),
		Visibility:                mergeVisibility(input.Visibility, existingRecord.Visibility),
		FunderWebsite:             mergeNullString(input.FunderWebsite, existingRecord.FunderWebsite),
		Description:               mergeNullString(input.Description, existingRecord.Description),
		EligibilityNotes:          mergeNullString(input.EligibilityNotes, existingRecord.EligibilityNotes),
		AwardAmountMin:            mergeFloat64(input.AwardAmountMin, existingRecord.AwardAmountMin),
		AwardAmountMax:            mergeFloat64(input.AwardAmountMax, existingRecord.AwardAmountMax),
		EstimatedApplicationHours: mergeFloat64(input.EstimatedApplicationHours, existingRecord.EstimatedApplicationHours),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return toGrantResponse(record), nil
}

func (s *GrantService) DeleteGrant(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID) error {
	err := s.repo.DeleteGrant(ctx, repository.DeleteGrantParams{
		GrantWriterID: grantWriterID,
		ID:            grantID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGrantNotFound
		}
		return err
	}

	return nil
}

func (s *GrantService) AddDeadline(ctx context.Context, input AddDeadlineInput) (*Deadline, error) {
	_, err := s.repo.GetGrantByID(ctx, repository.GetGrantByIDParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.GrantID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}

	if !isValidDeadlineLabel(input.Label) {
		return nil, errors.New("invalid deadline label")
	}

	record, err := s.repo.CreateDeadline(ctx, repository.CreateDeadlineParams{
		GrantID:     input.GrantID,
		Label:       repository.GrantDeadlineType(input.Label),
		Date:        input.Date,
		Description: toNullStringFromPtr(input.Description),
	})

	if err != nil {
		return nil, err
	}

	return toDeadlineResponse(record), nil
}

func (s *GrantService) DeleteDeadline(ctx context.Context, grantWriterID uuid.UUID, grantID uuid.UUID, deadlineID uuid.UUID) error {
	_, err := s.repo.GetGrantByID(ctx, repository.GetGrantByIDParams{
		GrantWriterID: grantWriterID,
		ID:            grantID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGrantNotFound
		}
		return err
	}

	err = s.repo.DeleteDeadline(ctx, repository.DeleteDeadlineParams{
		GrantWriterID: grantWriterID,
		GrantID:       grantID,
		ID:            deadlineID,
	})

	if err != nil {
		return err
	}

	return nil
}

func toGrantResponse(g repository.Grant) *Grant {
	return &Grant{
		ID:                        g.ID,
		GrantWriterID:             g.GrantWriterID,
		Title:                     g.Title,
		FunderName:                g.FunderName,
		Visibility:                string(g.Visibility),
		FunderWebsite:             g.FunderWebsite.String,
		Description:               g.Description.String,
		AwardAmountMin:            nullStringToFloat64(g.AwardAmountMin),
		AwardAmountMax:            nullStringToFloat64(g.AwardAmountMax),
		EligibilityNotes:          g.EligibilityNotes.String,
		EstimatedApplicationHours: nullStringToFloat64(g.EstimatedApplicationHours),
		CreatedAt:                 g.CreatedAt,
		UpdatedAt:                 g.UpdatedAt,
	}
}

func toDeadlineResponse(d repository.GrantDeadline) *Deadline {
	return &Deadline{
		ID:          d.ID,
		GrantID:     d.GrantID,
		Label:       string(d.Label),
		Date:        d.Date,
		Description: d.Description.String,
		CreatedAt:   d.CreatedAt,
	}
}

func isValidDeadlineLabel(label string) bool {
	switch repository.GrantDeadlineType(label) {
	case repository.GrantDeadlineTypeApplication,
		repository.GrantDeadlineTypeLOI,
		repository.GrantDeadlineTypeOther:
		return true
	}

	return false
}
