package util

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type config struct {
	Home         string         `env:"HOME"`
	Port         int            `env:"PORT" envDefault:"3000"`
	Password     string         `env:"PASSWORD,unset"`
	IsProduction bool           `env:"PRODUCTION"`
	Duration     time.Duration  `env:"DURATION"`
	Hosts        []string       `env:"HOSTS" envSeparator:":"`
	TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/tmp"`
	StringInts   map[string]int `env:"MAP_STRING_INT"`
}

func get_env() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

  // or you can use generics
  cfg, err := env.ParseAs[config]()
  if err != nil {
		fmt.Printf("%+v\n", err)
  }

	fmt.Printf("%+v\n", cfg)
}
