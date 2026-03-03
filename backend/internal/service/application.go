package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
)

var (
	ErrApplicationNotFound         = errors.New("application not found")
	ErrApplicationAlreadyPublished = errors.New("application already published")
)

type ApplicationService struct {
	repo repository.Querier
}

func NewApplicationService(repo repository.Querier) *ApplicationService {
	return &ApplicationService{repo: repo}
}

type Application struct {
	ID            uuid.UUID  `json:"id"`
	GrantWriterID uuid.UUID  `json:"grant_writer_id"`
	GrantID       uuid.UUID  `json:"grant_id"`
	ClientID      uuid.UUID  `json:"client_id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	IsExclusive   bool       `json:"is_exclusive"`
	PublishedAt   *time.Time `json:"published_at"`
	Notes         *string    `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     time.Time  `json:"deleted_at"`
}

type CreateApplicationInput struct {
	GrantWriterID uuid.UUID
	GrantID       uuid.UUID
	ClientID      uuid.UUID
	Title         string
	Status        string
	IsExclusive   bool
	Notes         *string
}

type UpdateApplicationInput struct {
	GrantWriterID uuid.UUID
	ApplicationID uuid.UUID
	Title         string
	Status        string
	IsExclusive   bool
	Notes         *string
}

func (s *ApplicationService) CreateApplication(ctx context.Context, input CreateApplicationInput) (*Application, error) {
	record, err := s.repo.CreateApplication(ctx, repository.CreateApplicationParams{
		GrantWriterID: input.GrantWriterID,
		GrantID:       input.GrantID,
		ClientID:      input.ClientID,
		Title:         input.Title,
		Status:        repository.ApplicationStatusNotStarted,
		IsExclusive:   input.IsExclusive,
		Notes:         toNullStringFromPtr(input.Notes),
	})
	if err != nil {
		return nil, err
	}

	return toApplicationResponse(record), nil

}

func (s *ApplicationService) GetAllApplicationsByUserID(ctx context.Context, grantWriterID uuid.UUID) ([]Application, error) {
	records, err := s.repo.GetAllApplicationsByUserID(ctx, grantWriterID)
	if err != nil {
		return nil, err
	}

	applications := make([]Application, len(records))
	for i, record := range records {
		applications[i] = *toApplicationResponse(record)
	}

	return applications, nil
}

func (s *ApplicationService) GetAllApplicationsByClientID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) ([]Application, error) {
	records, err := s.repo.GetAllApplicationsByClientID(ctx, repository.GetAllApplicationsByClientIDParams{
		GrantWriterID: grantWriterID,
		ClientID:      clientID,
	})
	if err != nil {
		return nil, err
	}

	applications := make([]Application, len(records))
	for i, record := range records {
		applications[i] = *toApplicationResponse(record)
	}

	return applications, nil
}

func (s *ApplicationService) GetApplicationByID(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) (*Application, error) {
	record, err := s.repo.GetApplicationByID(ctx, repository.GetApplicationByIDParams{
		GrantWriterID: grantWriterID,
		ID:            applicationID,
	})
	if err != nil {
		return nil, err
	}

	return toApplicationResponse(record), nil
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, input UpdateApplicationInput) (*Application, error) {
	existingRecord, err := s.repo.GetApplicationByID(ctx, repository.GetApplicationByIDParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.ApplicationID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, err
	}

	record, err := s.repo.UpdateApplication(ctx, repository.UpdateApplicationParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.ApplicationID,
		Title:         input.Title,
		Status:        repository.ApplicationStatus(input.Status),
		IsExclusive:   input.IsExclusive,
		Notes:         toNullString(mergeString(input.Notes, existingRecord.Notes.String)),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, err
	}

	return toApplicationResponse(record), nil

}

func (s *ApplicationService) PublishApplication(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) (*Application, error) {
	_, err := s.repo.GetApplicationByID(ctx, repository.GetApplicationByIDParams{
		GrantWriterID: grantWriterID,
		ID:            applicationID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}
		return nil, err
	}

	record, err := s.repo.PublishApplication(ctx, repository.PublishApplicationParams{
		GrantWriterID: grantWriterID,
		ID:            applicationID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationAlreadyPublished
		}
		return nil, err
	}

	return toApplicationResponse(record), nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, grantWriterID uuid.UUID, applicationID uuid.UUID) error {
	err := s.repo.DeleteApplication(ctx, repository.DeleteApplicationParams{
		GrantWriterID: grantWriterID,
		ID:            applicationID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrClientNotFound
		}
		return err
	}

	return nil
}

func toApplicationResponse(a repository.Application) *Application {
	return &Application{
		ID:            a.ID,
		GrantWriterID: a.GrantWriterID,
		GrantID:       a.GrantID,
		ClientID:      a.ClientID,
		Title:         a.Title,
		Status:        string(a.Status),
		IsExclusive:   a.IsExclusive,
		PublishedAt:   nullDateToTime(a.PublishedAt),
		Notes:         &a.Notes.String,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}
