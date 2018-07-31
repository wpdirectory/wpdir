package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config contains useful App data.
type Config struct {
	Version       string
	Name          string
	Commit        string
	Date          string
	WD            string
	UpdateWorkers int
	SearchWorkers int
	DNS           struct {
		Email  string
		APIKey string
	}
	Host  string
	Ports struct {
		HTTP  string
		HTTPS string
	}
}

// Setup creates and returns the Config struct.
func Setup() *Config {

	viper.SetDefault("version", "0.1.0")
	viper.SetDefault("name", "WP Directory")
	viper.SetDefault("commit", "")
	viper.SetDefault("date", "")
	viper.SetDefault("updateworkers", 4)
	viper.SetDefault("searchworkers", 6)
	viper.SetDefault("host", "http://localhost")
	viper.SetDefault("ports.http", "80")
	viper.SetDefault("ports.https", "443")

	viper.AddConfigPath("/etc/wpdir/")
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %s\n", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting Working Directory: %s\n", err)
	}

	config := &Config{
		Version:       "0.5.0",
		Name:          viper.GetString("name"),
		Commit:        "dgfsdghfisdfdsfsdf",
		Date:          time.Now().Format(time.RFC3339Nano),
		WD:            wd,
		UpdateWorkers: viper.GetInt("updateworkers"),
		SearchWorkers: viper.GetInt("searchworkers"),
		Host:          viper.GetString("host"),
	}

	config.Ports.HTTP = viper.GetString("ports.http")
	config.Ports.HTTPS = viper.GetString("ports.https")

	return config

}
