package source

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/saipanno/go-kit/logger"
	"github.com/saipanno/go-kit/utils"
)

type UnSplashItem struct {
	ID   string `json:"id,omitempty"`
	URLs struct {
		Raw string `json:"raw,omitempty"`
	} `json:"urls,omitempty"`
	Links struct {
		HTML string `json:"html,omitempty"`
	} `json:"links,omitempty"`
	Likes int `json:"likes,omitempty"`
}

func (ui *UnSplashItem) ToItem() (item *LikeItem) {

	now := time.Now()

	item = &LikeItem{
		Source:   "unsplash.com",
		SourceID: ui.ID,
		Refer:    ui.Links.HTML,
		RawURL:   ui.URLs.Raw,
		FetchAt:  &now,
	}
	item.BuildID()
	return
}

type UnSplash struct {
	Username  string
	AccessKey string
	DataDir   string
}

func (us *UnSplash) SetConfig(sc map[string]string) error {

	v, has := sc["username"]
	if !has {
		return errors.New("username field is need")
	}
	us.Username = v

	v, has = sc["access_key"]
	if !has {
		return errors.New("access_key field is need")
	}
	us.AccessKey = v

	v, has = sc["data_dir"]
	if !has {
		return errors.New("data_dir field is need")
	}
	us.DataDir = v

	return nil
}

func (us *UnSplash) Fetch() (data []*LikeItem, err error) {

	var page = 1
	var pageSize = 25
	var header = make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Client-ID %s", us.AccessKey))

	option := utils.NewOptions(
		utils.WithHeader(header),
	)

	for {

		var likes []*UnSplashItem

		err = utils.GetURLWithJSONResult(
			fmt.Sprintf("https://api.unsplash.com/users/%s/likes?page=%d&per_page=%d", us.Username, page, pageSize),
			&likes,
			option,
		)
		if err != nil {
			return
		}

		for _, item := range likes {
			data = append(data, item.ToItem())
		}

		if len(likes) < pageSize {
			break
		}

		page++
	}

	logger.Infof("total %d like item on unsplash.com", len(data))
	return
}

func (us *UnSplash) Download(item *LikeItem) (err error) {

	return utils.DownloadFile(item.RawURL, filepath.Join(us.DataDir, item.ID))
}
