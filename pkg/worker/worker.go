package worker

import (
	"ShoppingList-Backend/pkg/application"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Nerzal/gocloak/v8"
	"github.com/gocraft/work"
	"go.uber.org/zap"
)

const Namespace = "ShoppingListV4"

type WorkerContext struct {
	App *application.Application
}

func NewWorkerPool(app *application.Application) *work.WorkerPool {
	pool := work.NewWorkerPool(WorkerContext{}, 10, Namespace, app.Redis)
	pool.Middleware((*WorkerContext).Log)
	pool.Middleware(func(c *WorkerContext, job *work.Job, next work.NextMiddlewareFunc) error {
		c.App = app
		return next()
	})
	pool.Job(JobDemoCleanUp, (*WorkerContext).CleanUpDemoUsers)

	return pool
}

func Start(pool *work.WorkerPool, wg *sync.WaitGroup) {
	defer wg.Done()
	pool.Start()
	zap.S().Info("Worker pool started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	pool.Stop()
}

func (c *WorkerContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	zap.S().Infow("Starting job", "job name", job.Name)
	return next()
}

func (c *WorkerContext) CleanUpDemoUsers(job *work.Job) error {
	client := gocloak.NewClient(c.App.Cfg.JwtKeycloakUrl)
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, c.App.Cfg.JwtKeycloakUsername, c.App.Cfg.JwtKeycloakPassword, "master")
	if err != nil {
		zap.S().Errorf("Could not login: %v", err)
		return err
	}

	realm := "shoppinglist"
	demoUsersGroupId := "eab4732c-525c-4456-926d-c88b8bc0a55a"

	// zap.S().Infow("Group", "demoUsersGroup", demoUsersGroup)

	demoUsers, err := client.GetGroupMembers(ctx, token.AccessToken, realm, demoUsersGroupId, gocloak.GetGroupsParams{})
	if err != nil {
		zap.S().Errorf("Could not get users of demo group (id: %v): %v", err, demoUsersGroupId)
		return err
	}
	// zap.S().Infow("Users", "users", demoUsers)
	for _, user := range demoUsers {
		userID := user.ID
		// zap.S().Infow("user", "userID", user.ID)
		if err := c.App.Queries.Item.DeleteItems(*userID); err != nil {
			zap.S().Errorf("Error deleting items: %v", err)
			return err
		}
		if err := c.App.Queries.List.DeleteLists(*userID); err != nil {
			zap.S().Errorf("Error deleting lists: %v", err)
			return err
		}
	}

	zap.S().Infow("Finished job", "job name", job.Name)
	return nil
}
