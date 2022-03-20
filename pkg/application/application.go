package application

import (
	"ShoppingList-Backend/internal/pkg/item"
	"ShoppingList-Backend/internal/pkg/list"
	"ShoppingList-Backend/pkg/config"
	"ShoppingList-Backend/pkg/db"
	"ShoppingList-Backend/pkg/server"
	"ShoppingList-Backend/pkg/sse"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Application struct {
	Cfg                   *config.Config
	Queries               *Repositories
	Controllers           *Controllers
	Redis                 *redis.Pool
	Srv                   *server.Server
	SseBroker             *sse.Broker
	GetRabbitMqConnection func() (*amqp.Connection, error)
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

	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.GetRedisConnStr(), redis.DialClientName(cfg.GetRedisClientName()), redis.DialUsername(cfg.GetRedisUser()), redis.DialPassword(cfg.GetRedisPassword()))
		},
	}

	rabbitMqConnectionGetter := func() func() (*amqp.Connection, error) {

		connectToRabbitMQ := func() (*amqp.Connection, error) {
			hostname, err := os.Hostname()
			if err != nil {
				zap.S().Error("could not get os.Hostname: %w", err)
			}
			conn, err := amqp.DialConfig(cfg.GetRabbitMqUri(), amqp.Config{
				Properties: amqp.Table{
					"connection_name": hostname + "-ShoppingList.API",
				},
			})
			return conn, err
		}

		conn, err := connectToRabbitMQ()
		return func() (*amqp.Connection, error) {
			if conn.IsClosed() {
				conn, err = connectToRabbitMQ()
			}
			return conn, err
		}
	}

	getRabbitMqConn := rabbitMqConnectionGetter()

	sseBroker := sse.NewBroker(getRabbitMqConn)

	return &Application{
		Cfg:                   cfg,
		Queries:               repos,
		Redis:                 redisPool,
		Controllers:           controllers,
		SseBroker:             sseBroker,
		GetRabbitMqConnection: getRabbitMqConn,
	}, nil
}
