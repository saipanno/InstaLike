package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/saipanno/go-kit/client"
	"github.com/saipanno/go-kit/logger"
	"github.com/saipanno/go-kit/utils"
)

var (
	ac     *AppConfig
	locker sync.RWMutex
)

type SourceConfig struct {
	AccessKey string `json:"access_key,omitempty"`
	DataDir   string `json:"data_dir,omitempty"`
}

// AppConfig ...
type AppConfig struct {
	Mode    string                       `json:"mode,omitempty"`
	Logger  *logger.Config               `json:"logger,omitempty"`
	DB      *client.DBConfig             `json:"db,omitempty"`
	Sources map[string]map[string]string `json:"sources,omitempty"`
}

// ParseConfig ...
func ParseConfig(f string) (err error) {

	var buf []byte
	var cfg AppConfig

	if !utils.FileExist(f) {
		err = fmt.Errorf("Config file(%s) not exists", f)
		return
	}

	buf, err = ioutil.ReadFile(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return
	}

	locker.Lock()
	defer locker.Unlock()
	ac = &cfg

	logger.Debugf("config is %s", utils.PrettyPrint(ac))
	return
}

// Config ...
func Config() *AppConfig {

	locker.RLock()
	defer locker.RUnlock()

	return ac
}
