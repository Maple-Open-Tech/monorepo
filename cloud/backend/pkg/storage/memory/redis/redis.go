package redis

import (
	"context"
	"errors"
	"log"
	"time"

	"log/slog"

	"github.com/redis/go-redis/v9"

	c "github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
)

type Cacher interface {
	Shutdown(ctx context.Context)
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte) error
	SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
	Publish(ctx context.Context, channelName string, binary []byte) error
	Subscribe(ctx context.Context, channelName string) RedisSubscriber
	GetRedisClient() redis.UniversalClient
}

type cache struct {
	Client *redis.Client
	Logger *slog.Logger
}

func NewCache(cfg *c.Configuration, logger *slog.Logger) Cacher {
	logger.Debug("cache initializing...")

	opt, err := redis.ParseURL(cfg.Cache.URI)
	if err != nil {
		logger.Error("cache failed parsing url", slog.Any("err", err), slog.String("URI", cfg.Cache.URI))
		log.Fatal(err)
	}
	rdb := redis.NewClient(opt)

	// Confirm connection with Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5-second timeout for initialization
	defer cancel()

	// Confirm connection made with Redis
	for i := 0; i < 3; i++ { // Retry logic
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			break
		}
		if i == 2 {
			logger.Error("cache failed connecting to Redis", slog.Any("err", err), slog.String("URI", cfg.Cache.URI))
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}

	logger.Debug("cache initialized successfully")
	return &cache{
		Client: rdb,
		Logger: logger,
	}
}

func (s *cache) GetRedisClient() redis.UniversalClient {
	return s.Client
}

func (s *cache) Shutdown(ctx context.Context) {
	s.Logger.Debug("cache shutting down...")

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.Client.Close()
	}()

	select {
	case <-done:
		s.Logger.Debug("cache shutdown complete")
	case <-ctx.Done():
		s.Logger.Warn("cache shutdown timeout or context canceled")
	}
}

func (s *cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := s.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // Key does not exist
	}
	if err != nil {
		s.Logger.Error("cache get failed", slog.Any("error", err))
		return nil, err
	}
	return []byte(val), nil
}

func (s *cache) Set(ctx context.Context, key string, val []byte) error {
	err := s.Client.Set(ctx, key, val, 0).Err()
	if err != nil {
		s.Logger.Error("cache set failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error {
	err := s.Client.Set(ctx, key, val, expiry).Err()
	if err != nil {
		s.Logger.Error("cache set with expiry failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) Delete(ctx context.Context, key string) error {
	err := s.Client.Del(ctx, key).Err()
	if err != nil {
		s.Logger.Error("cache delete failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) Publish(ctx context.Context, channelName string, binary []byte) error {
	err := s.Client.Publish(ctx, channelName, binary).Err()
	if err != nil {
		s.Logger.Error("cache failed publishing", slog.Any("error", err), slog.String("channel", channelName))
		return err
	}
	return nil
}

type RedisSubscriber interface {
	WaitUntilReceiveMessage(ctx context.Context) ([]byte, error)
	Close() error
}

func (s *cache) Subscribe(ctx context.Context, channelName string) RedisSubscriber {
	pubsub := s.Client.Subscribe(ctx, channelName)
	return &redisSubscriberImpl{
		pubsub: pubsub,
		logger: s.Logger,
	}
}

type redisSubscriberImpl struct {
	pubsub *redis.PubSub
	logger *slog.Logger
}

func (s *redisSubscriberImpl) WaitUntilReceiveMessage(ctx context.Context) ([]byte, error) {
	msg, err := s.pubsub.ReceiveMessage(ctx)
	if err != nil {
		s.logger.Error("failed to receive message", slog.Any("error", err))
		return nil, err
	}
	return []byte(msg.Payload), nil
}

func (s *redisSubscriberImpl) Close() error {
	return s.pubsub.Close()
}
