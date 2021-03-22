package ilike

import (
	"context"
	"time"

	"github.com/saipanno/InstaLike/source"
	"github.com/saipanno/go-kit/logger"
)

func (manager *Manager) CreateLike(item *source.LikeItem) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = manager.db.NamedExecContext(
		ctx,
		`INSERT INTO likes
				(id, source, source_id, refer, raw_url, fetch_at)
		VALUES  
				(:id, :source, :source_id, :refer, :raw_url, :fetch_at);`,
		item)
	if err != nil {
		logger.Errorf("insert like item failed, message is %s", err.Error())
		return
	}

	manager.likes.Store(item.ID, item)

	logger.Debugf("item(%s|%s) create success", item.Source, item.SourceID)
	return
}

func (manager *Manager) FetchLikes() (data []*source.LikeItem, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = manager.db.SelectContext(ctx, &data, `SELECT * FROM likes;`)
	if err != nil {
		logger.Errorf("select like item failed, message is %s", err.Error())
		return
	}

	logger.Infof("total %d like item in db", len(data))
	return
}
