package store

import (
	"context"
	"fmt"
	"time"

	"github.com/gitwub5/go_todo_app/config"
	"github.com/gitwub5/go_todo_app/entity"
	"github.com/go-redis/redis/v8"
)

/*
Redis를 사용해 액세스 토큰을 관리하는 키-값 저장소
*/

// 새로운 KVS를 생성한다.
func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

// KVS는 키-값 저장소를 나타낸다.
type KVS struct {
	Cli *redis.Client
}

// Save는 주어진 키에 주어진 사용자 ID를 저장한다.
func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	id := int64(userID)
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

// Load는 주어진 키에 저장된 사용자 ID를 반환한다.
func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
