package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// CF -> for use configs model
	CF = &Configs{}
)

// Environment environment
type Environment string

const (
	// DEV environment develop
	DEV Environment = "dev"
	// PROD environment production
	PROD Environment = "prod"
)

// Configs models
type Configs struct {
	Validator *validator.Validate
	Database  struct {
		MySQL struct {
			Host         string `mapstructure:"host"`
			Port         string `mapstructure:"port"`
			Username     string `mapstructure:"username"`
			Password     string `mapstructure:"password"`
			DatabaseName string `mapstructure:"database_name"`
			DriverName   string `mapstructure:"driver_name"`
			Protocol     string `mapstructure:"protocol"`
			Timeout      string `mapstructure:"timeout"`
		} `mapstructure:"mysql"`
	} `mapstructure:"database"`
	Swagger struct {
		Title       string   `mapstructure:"title"`
		Description string   `mapstructure:"description"`
		Version     string   `mapstructure:"version"`
		Host        string   `mapstructure:"host"`
		BasePath    string   `mapstructure:"base_path"`
		Schemes     []string `mapstructure:"schemes"`
	} `mapstructure:"swagger"`
	HTTPServer struct {
		ReadTimeout       time.Duration `mapstructure:"read_timeout"`
		WriteTimeout      time.Duration `mapstructure:"write_timeout"`
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	} `mapstructure:"http_server"`
}

// InitConfig init config
func InitConfig(configPath string, environment string) error {
	v := viper.New()

	v.AddConfigPath(configPath)
	v.SetConfigName(fmt.Sprintf("config.%s", CF.parseEnvironment(environment)))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := bindingConfig(v, CF); err != nil {
			logrus.Error("binding error:", err)
			return
		}
	})

	return nil
}

// bindingConfig binding config
func bindingConfig(vp *viper.Viper, cf *Configs) error {
	if err := vp.Unmarshal(&cf); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}

	cf.Validator = validator.New()

	return nil
}

func (c Configs) parseEnvironment(environment string) Environment {
	switch environment {
	case "dev":
		return DEV

	case "prod":
		return PROD
	}

	return DEV
}
