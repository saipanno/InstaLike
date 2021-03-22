package source

import (
	"time"

	"github.com/saipanno/go-kit/utils"
)

type LikeItem struct {
	ID       string     `db:"id"`
	Source   string     `db:"source"`
	SourceID string     `db:"source_id"`
	Refer    string     `db:"refer"`
	RawURL   string     `db:"raw_url"`
	FetchAt  *time.Time `db:"fetch_at"`
}

func (item *LikeItem) BuildID() {
	item.ID = utils.Md5sum(item.Source + item.SourceID)
}

type Plugin interface {
	Fetch() ([]*LikeItem, error)
	Download(*LikeItem) error
}
