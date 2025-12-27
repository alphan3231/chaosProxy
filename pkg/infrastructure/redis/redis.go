package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr, password string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

// GetRawClient returns the underlying redis client if needed
func (c *Client) GetRawClient() *redis.Client {
	return c.rdb
}

type TrafficLog struct {
	Timestamp    time.Time `json:"timestamp"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	RequestBody  string    `json:"request_body"`
	Status       int       `json:"status"`
	ResponseBody string    `json:"response_body"`
	Duration     string    `json:"duration"`
}

func (c *Client) PublishTraffic(ctx context.Context, logEntry TrafficLog) error {
	// For now, we just JSON encode and publish to a channel
	// In a real high-perf scenario, we might use a worker pool/buffer here
	return c.rdb.Publish(ctx, "chaos:traffic", logEntry).Err()
}

func (c *Client) GetGhostResponse(ctx context.Context, method, path string) (*TrafficLog, error) {
	key := fmt.Sprintf("chaos:ghost:%s:%s", method, path)
	data, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var logEntry TrafficLog
	if err := json.Unmarshal([]byte(data), &logEntry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ghost response: %w", err)
	}
	return &logEntry, nil
}
