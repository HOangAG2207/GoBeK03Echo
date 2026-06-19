package infrastructure

import (
	"log"

	"github.com/HOangAG2207/GoBeK03Echo/internal/api"
	pkglogger "github.com/HOangAG2207/GoBeK03Echo/pkg/logger"
	pkgredis "github.com/HOangAG2207/GoBeK03Echo/pkg/redis"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
func CreateAPI() api.Engine {
	cfg := CreateAPIConfig()

	redisClient, err := pkgredis.NewClient("")
	if err != nil {
		panic(err)
	}
	if err := pkglogger.SetLogLevel(); err != nil {
		panic(err)
	}
	appEngine := api.NewEngine(&api.EngineOpts{
		Cfg:           cfg,
		Redis:         redisClient,
		RandomCodeGen: pkgutils.NewCodeGenerator(),
	})
	return appEngine
}
