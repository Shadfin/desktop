package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/mleku/appdata"
	"github.com/vmihailenco/msgpack/v5"
)

var ()

type Config struct {
	ctx  context.Context `json:"-"`
	path string
	data *ConfigData
	mu   *sync.Mutex
}

type ConfigData struct {
	Authentication *AuthenticationStore `json:"authentication"`
	Server         *ServerStore         `json:"server"`
	User           *UserStore           `json:"user"`
}

// NewApp creates a new App application struct
func NewConfig() *Config {
	cfg := &Config{mu: &sync.Mutex{}}
	homeDir := UserHomeDir()
	println("Home Directory is " + homeDir)
	path_cf := appdata.GetDataDir("com.shadfin.app", true, "")

	cfg.path = path.Join(homeDir + "\\com.shadfin.app\\shadfin_app.cfg")

	fmt.Printf("Config dir is %v\n", path_cf)

	cfg.WriteIfNotExists()
	cfg.Read()
	return cfg
}

func (c *Config) GetData() *ConfigData {
	return c.data
}

func (c *Config) SetAuthentication(store *AuthenticationStore) {
	c.data.Authentication = store
	c.Write()
}

func (c *Config) SetServer(store *ServerStore) {
	c.data.Server = store
	c.Write()
}

func (c *Config) SetUser(store *UserStore) {
	c.data.User = store
	c.Write()
}

func (c *Config) WriteIfNotExists() {
	if _, err := os.Stat(c.path); errors.Is(err, os.ErrNotExist) {
		println("Config file did not exist, writing one now...")
		c.Write()
	}
}
func (c *Config) Read() {
	c.mu.Lock()
	file, err := ioutil.ReadFile(c.path)
	c.mu.Unlock()

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	cfg := ConfigData{}

	err = msgpack.Unmarshal(file, &cfg)

	if err != nil {
		log.Fatal("Error when unpacking: ", err)
	}

	c.data = &cfg
}

func (c *Config) Write() {

	cfg_bytes, err := msgpack.Marshal(c.data)

	if err != nil {
		panic(err)
	}
	println("Trying to write: ", string(cfg_bytes))
	c.mu.Lock()
	err = os.WriteFile(c.path, cfg_bytes, 0777)
	c.mu.Unlock()

	if err != nil {
		panic(err)
	}

	println("Successfully wrote file")
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (c *Config) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("APPDATA")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
