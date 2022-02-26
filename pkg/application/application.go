package application

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/internal/pkg/list"
	"ShoppingList-Backend/internal/pkg/queries"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/db"

	"github.com/gomodule/redigo/redis"
)

type Application struct {
	Cfg     *config.Config
	Queries *queries.Queries
	Redis   *redis.Pool
}

func Get(cfg *config.Config) (*Application, error) {
	db, err := db.Get(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	queries := &queries.Queries{
		Item: &item.ItemQueries{
			DB: db.Client,
		},
		List: &list.ListQueries{
			DB: db.Client,
		},
	}

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.GetRedisConnStr(), redis.DialClientName(cfg.GetRedisClientName()), redis.DialUsername(cfg.GetRedisUser()), redis.DialPassword(cfg.GetRedisPassword()))
		},
	}

	return &Application{
		Cfg:     cfg,
		Queries: queries,
		Redis:   redisPool,
	}, nil
}
