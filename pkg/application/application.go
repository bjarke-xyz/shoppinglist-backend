package application

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/internal/pkg/list"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/db"
	"ShoppingList-Backend/pkg/server"

	"github.com/gomodule/redigo/redis"
)

type Application struct {
	Cfg         *config.Config
	Queries     *Repositories
	Controllers *Controllers
	Redis       *redis.Pool
	Srv         *server.Server
}

func Get(cfg *config.Config) (*Application, error) {
	db, err := db.Get(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	repos := &Repositories{
		Item: &item.ItemRepository{
			DB: db.Client,
		},
		List: &list.ListRepository{
			DB: db.Client,
		},
	}

	controllers := &Controllers{
		Item: item.NewItemController(repos.Item),
		List: list.NewListController(repos.Item, repos.List),
	}

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.GetRedisConnStr())
		},
	}

	return &Application{
		Cfg:         cfg,
		Queries:     repos,
		Redis:       redisPool,
		Controllers: controllers,
	}, nil
}
