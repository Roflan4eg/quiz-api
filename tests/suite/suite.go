package suite

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Roflan4eg/quiz-api/config"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	BaseURL    string
	Client     *http.Client
	TestUserID string
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	baseURL := "http://localhost:" + cfg.HTTP.Port
	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		BaseURL:    baseURL,
		Client:     client,
		TestUserID: "test-user-" + time.Now().Format("20060102150405"),
	}
}

func loadConfig() (*config.Config, error) {
	var cfg config.Config
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
