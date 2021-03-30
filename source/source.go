package source

import (
	"fmt"
	"sync"
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

type Source interface {
	SetConfig(map[string]string) error
	Fetch() ([]*LikeItem, error)
	Download(*LikeItem) error
}

var (
	locker  sync.RWMutex
	sources = make(map[string]Source)
)

func init() {
	register("unsplash", &UnSplash{})
}

func register(name string, src Source) {

	locker.Lock()
	defer locker.Unlock()

	if src == nil {
		panic("source is nil")
	}

	if _, exist := sources[name]; exist {
		panic("source already registered")
	}

	sources[name] = src
}

// New ...
func New(name string, sc map[string]string) (src Source, err error) {

	locker.RLock()
	defer locker.RUnlock()

	var exist bool
	if src, exist = sources[name]; exist {

		err = src.SetConfig(sc)
		if err != nil {
			return
		}

		return
	}

	err = fmt.Errorf("source(%s) is not exist", name)
	return
}
