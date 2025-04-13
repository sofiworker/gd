package conf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/chuck1024/gd/v2/logger"
	"github.com/spf13/viper"
)

var (
	v atomic.Value
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", ""))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", ""))
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("not found config in ['/etc', '$HOME/', './conf/']")
		} else {
			logger.Fatalf("read config failed:%v", err)
		}
	}
	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		logger.Errorf("unmarshal config to struct failed:%v", err)
	}
	v.Store(&conf)
}

type Conf interface {
	SetSource(rw io.ReadWriter) error
	ReadRaw() error
	Read() error
	Get() (interface{}, error)
	Bind(v interface{}) error
	Set() error
	Watch()
}

type ConfigOptsFunc func(*ConfigOpts)

type ConfigOpts struct {
	Name  string
	Ext   string
	Paths []string
}

func ConfigDefaultPaths() []string {
	userHomeDir, _ := os.UserHomeDir()
	return []string{userHomeDir}
}

func WithConfigPaths(paths []string) ConfigOptsFunc {
	return func(opts *ConfigOpts) {
		opts.Paths = paths
	}
}

func WithConfigName(name string) ConfigOptsFunc {
	return func(opts *ConfigOpts) {
		opts.Name = name
	}
}

func WithConfigExt(ext string) ConfigOptsFunc {
	return func(opts *ConfigOpts) {
		opts.Ext = ext
	}
}

type Config struct {
	opts *ConfigOpts
}

func NewConfig(opts ...ConfigOptsFunc) (Conf, error) {
	var configOpts ConfigOpts
	for _, opt := range opts {
		opt(&configOpts)
	}
	return &Config{
		opts: &configOpts,
	}, nil
}

func (c *Config) ReadRaw() error {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Read() error {
	return c.ReadRaw()
}

func (c *Config) Get() (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Bind(v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Set() error {
	//TODO implement me
	panic("implement me")
}

func (c *Config) Watch() {
	//TODO implement me
	panic("implement me")
}

func (c *Config) SetSource(rw io.ReadWriter) error {
	return nil
}


func DefaultSource() io.ReadWriter {
	return bytes.NewBuffer(nil)
}
