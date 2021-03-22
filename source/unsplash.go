package source

import (
	"fmt"
	"net/http"
	"time"

	"github.com/saipanno/go-kit/logger"
	"github.com/saipanno/go-kit/utils"
)

var (
	unsplashAccessKey = "RCslKktjieC-C6DNWeglKwnxfQP0P4ogPrEBvL4nVeg"
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
}

func (us *UnSplash) Fetch() (data []*LikeItem, err error) {

	var page = 1
	var pageSize = 25
	var header = make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Client-ID %s", unsplashAccessKey))

	option := utils.NewOptions(
		utils.WithHeader(header),
	)

	for {

		var likes []*UnSplashItem

		err = utils.GetURLWithJSONResult(
			fmt.Sprintf("https://api.unsplash.com/users/saipanno/likes?page=%d&per_page=%d", page, pageSize),
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

	return utils.DownloadFile(item.RawURL, fmt.Sprintf("../../data/%s.jpg", item.ID))
}
