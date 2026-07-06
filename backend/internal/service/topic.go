package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jbechler2/grant-tool/backend/internal/repository"
)

var (
	ErrTopicNotFound = errors.New("Topic not found")
)

type TopicService struct {
	repo repository.Querier
}

func NewTopicService(repo repository.Querier) *TopicService {
	return &TopicService{repo: repo}
}

type Topic struct {
	ID            uuid.UUID `json:"id"`
	GrantWriterID uuid.UUID `json:"grant_writer_id"`
	Label         string    `json:"label"`
}

type CreateTopicInput struct {
	GrantWriterID uuid.UUID
	Label         string
}

func (s *TopicService) CreateTopic(ctx context.Context, input CreateTopicInput) (*Topic, error) {
	record, err := s.repo.CreateTopic(ctx, repository.CreateTopicParams{
		GrantWriterID: input.GrantWriterID,
		Label:         input.Label,
	})
	if err != nil {
		return nil, err
	}

	return toTopicResponseFromRepositoryTopic(record), nil
}

func (s *TopicService) DeleteTopic(ctx context.Context, grantWriterID uuid.UUID, topicID uuid.UUID) error {
	err := s.repo.DeleteTopic(ctx, repository.DeleteTopicParams{
		GrantWriterID: grantWriterID,
		ID:            topicID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGrantNotFound
		}
		return err
	}

	return nil
}

func toTopicResponseFromRepositoryTopic(topic repository.Topic) *Topic {
	return &Topic{
		ID:            topic.ID,
		Label:         topic.Label,
		GrantWriterID: topic.GrantWriterID,
	}
}
