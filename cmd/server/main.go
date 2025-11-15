package main

import (
	"fmt"

	"github.com/Roflan4eg/quiz-api/config"
	"github.com/Roflan4eg/quiz-api/internal/app"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}
	logger.Init(cfg)

	newApp := app.New(cfg)

	if err = newApp.Setup(); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	if err = newApp.Start(); err != nil {
		logger.Error(err.Error())
	}

}
