package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Storage struct {
	client *redis.Client
}

type itemType string

const (
	itemFile itemType = "file"
	itemText itemType = "text"
)

type sharedItem struct {
	Type  itemType `json:"type"`
	Value string   `json:"value"`
}

func (s *sharedItem) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *sharedItem) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Add(ctx context.Context, item *sharedItem) error {
	key := strings.ReplaceAll(uuid.New().String(), "-", "")

	err := s.client.Set(ctx, key, item, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListAll(ctx context.Context) ([]*sharedItem, error) {
	items := make([]*sharedItem, 0)

	keys, err := s.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		item := &sharedItem{}

		op := s.client.Get(ctx, key)

		err := op.Err()
		if err != nil {
			return nil, err
		}

		err = op.Scan(item)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
