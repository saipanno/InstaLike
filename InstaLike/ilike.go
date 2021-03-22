package ilike

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	cronv3 "github.com/robfig/cron/v3"
	"github.com/saipanno/InstaLike/source"
	"github.com/saipanno/go-kit/logger"
	"github.com/saipanno/go-kit/utils"
)

type Manager struct {
	db        *sqlx.DB
	scheduler *cronv3.Cron

	likes   sync.Map
	sources []source.Plugin

	wg   sync.WaitGroup
	exit chan struct{}
}

// NewManager ...
func NewManager() *Manager {

	manager := &Manager{
		scheduler: cronv3.New(),
		exit:      make(chan struct{}),
	}

	return manager
}

func (manager *Manager) Start() (err error) {

	// Init Source
	manager.sources = append(manager.sources, &source.UnSplash{})

	// Init DB
	manager.db, err = sqlx.Connect("mysql", "root:252020@tcp(127.0.0.1:3306)/InstaLike?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logger.Errorf("create mysql connect failed, message is %s", err.Error())
		return
	}

	// Load Cache
	var likes []*source.LikeItem
	likes, err = manager.FetchLikes()
	if err != nil {
		return
	}

	for _, item := range likes {
		manager.likes.Store(item.ID, item)
	}

	// Init Cron
	manager.scheduler.AddFunc("*/2 * * * *", func() {

		var wg utils.WaitGroupWrapper
		for _, src := range manager.sources {
			wg.Wrap(func() {

				manager.wg.Add(1)
				defer manager.wg.Done()

				data, err1 := src.Fetch()
				if err1 != nil {
					return
				}

				for _, item := range data {

					if _, ok := manager.likes.Load(item.ID); ok {
						logger.Debugf("item(%s|%s) has extist", item.Source, item.SourceID)
						continue
					}

					logger.Debugf("item(%s) %s|%s has not extist", item.ID, item.Source, item.SourceID)

					err1 = src.Download(item)
					if err1 != nil {
						continue
					}

					err1 = manager.CreateLike(item)
					if err1 != nil {
						continue
					}

					time.Sleep(time.Second)
					logger.Debugf("item(%s|%s) download success", item.Source, item.ID)
				}
			})
		}

		wg.Wait()
		logger.Info("refresh loop finished")
	})
	manager.scheduler.Start()

	// Loop
	select {
	case <-manager.exit:
		return
	}
}

func (manager *Manager) Stop() (err error) {

	manager.scheduler.Stop()

	close(manager.exit)

	manager.wg.Wait()
	return
}
