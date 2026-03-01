package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
)

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrClientUnauthorized = errors.New("unauthorized access to client")
)

type ClientService struct {
	repo repository.Querier
}

func NewClientService(repo repository.Querier) *ClientService {
	return &ClientService{repo: repo}
}

type Client struct {
	ID            uuid.UUID `json:"id"`
	GrantWriterID uuid.UUID `json:"grant_writer_id"`
	Name          string    `json:"name"`
	ContactName   string    `json:"contact_name"`
	ContactPhone  string    `json:"contact_phone"`
	ContactEmail  string    `json:"contact_email"`
	Notes         string    `json:"notes"`
}

type CreateClientInput struct {
	GrantWriterID uuid.UUID
	Name          *string
	ContactName   *string
	ContactPhone  *string
	ContactEmail  *string
	Notes         *string
}

type UpdateClientInput struct {
	GrantWriterID uuid.UUID
	ID            uuid.UUID
	Name          *string
	ContactName   *string
	ContactPhone  *string
	ContactEmail  *string
	Notes         *string
}

func (s *ClientService) CreateClient(ctx context.Context, input CreateClientInput) (*Client, error) {
	record, err := s.repo.CreateClient(ctx, repository.CreateClientParams{
		Name:          *input.Name,
		GrantWriterID: input.GrantWriterID,
		ContactName:   toNullStringFromPtr(input.ContactEmail),
		ContactPhone:  toNullStringFromPtr(input.ContactPhone),
		ContactEmail:  toNullStringFromPtr(input.ContactEmail),
		Notes:         toNullStringFromPtr(input.Notes),
	})
	if err != nil {
		return nil, err
	}

	return toClientResponse(record), nil
}

func (s *ClientService) GetClientByID(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) (*Client, error) {
	record, err := s.repo.GetClientByID(ctx, repository.GetClientByIDParams{
		GrantWriterID: grantWriterID,
		ID:            clientID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return toClientResponse(record), nil
}

func (s *ClientService) GetAllClients(ctx context.Context, grantWriterID uuid.UUID) ([]Client, error) {
	records, err := s.repo.GetAllClientsByGrantWriter(ctx, grantWriterID)
	if err != nil {
		return nil, err
	}

	clients := make([]Client, len(records))
	for i, record := range records {
		clients[i] = *toClientResponse(record)
	}

	return clients, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, input UpdateClientInput) (*Client, error) {
	existingRecord, err := s.repo.GetClientByID(ctx, repository.GetClientByIDParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.ID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	record, err := s.repo.UpdateClient(ctx, repository.UpdateClientParams{
		GrantWriterID: input.GrantWriterID,
		ID:            input.ID,
		Name:          mergeString(input.Name, existingRecord.Name),
		ContactName:   mergeNullString(input.ContactName, existingRecord.ContactName),
		ContactPhone:  mergeNullString(input.ContactPhone, existingRecord.ContactPhone),
		ContactEmail:  mergeNullString(input.ContactEmail, existingRecord.ContactEmail),
		Notes:         mergeNullString(input.Notes, existingRecord.Notes),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return toClientResponse(record), nil
}

func (s *ClientService) DeleteClient(ctx context.Context, grantWriterID uuid.UUID, clientID uuid.UUID) error {
	err := s.repo.DeleteClient(ctx, repository.DeleteClientParams{
		GrantWriterID: grantWriterID,
		ID:            clientID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrClientNotFound
		}
		return err
	}

	return nil
}

func toClientResponse(c repository.Client) *Client {
	return &Client{
		ID:            c.ID,
		GrantWriterID: c.GrantWriterID,
		Name:          c.Name,
		ContactName:   c.ContactName.String,
		ContactPhone:  c.ContactPhone.String,
		ContactEmail:  c.ContactEmail.String,
		Notes:         c.Notes.String,
	}
}

func mergeNullString(input *string, existing sql.NullString) sql.NullString {
	if input == nil {
		return existing
	}

	return toNullStringFromPtr(input)
}
