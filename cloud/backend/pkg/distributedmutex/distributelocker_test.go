package distributedmutex

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// mockRedisClient implements minimal required methods
type mockRedisClient struct {
	redis.UniversalClient
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *mockRedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return redis.NewCmd(ctx)
}

func (m *mockRedisClient) EvalSha(ctx context.Context, sha string, keys []string, args ...interface{}) *redis.Cmd {
	return redis.NewCmd(ctx)
}

func (m *mockRedisClient) ScriptExists(ctx context.Context, scripts ...string) *redis.BoolSliceCmd {
	return redis.NewBoolSliceCmd(ctx)
}

func (m *mockRedisClient) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func TestNewAdapter(t *testing.T) {
	adapter := NewAdapter(slog.Default(), &mockRedisClient{})
	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
}

func TestAcquireAndRelease(t *testing.T) {
	ctx := context.Background()
	adapter := NewAdapter(slog.Default(), &mockRedisClient{})

	adapter.Acquire(ctx, "test-key")
	adapter.Acquiref(ctx, "test-key-%d", 1)
	adapter.Release(ctx, "test-key")
	adapter.Releasef(ctx, "test-key-%d", 1)
}
