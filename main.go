package main

import (
	"encoding/gob"
	"fmt"

	"git.betfavorit.cf/backend/logger"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"
	//"os/user"
	"net/http"
	"time"

	"os"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/application"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"
)

func init() {
	gob.Register(&models.UserRow{})
}

func newConfig() (*viper.Viper, error) {
	//u, err := user.Current()
	//if err != nil {
	//	return nil, err
	//}

	c := viper.New()

	c.SetEnvPrefix("PG")
	c.BindEnv("USER")
	pguser := c.Get("USER")
	c.BindEnv("PASSWORD")
	pgpass := c.Get("PASSWORD")
	c.BindEnv("DATABASE")
	pgdatabase := c.Get("DATABASE")

	c.SetEnvPrefix("KUBERNETES")
	c.BindEnv("ADDR")
	kubAddress := c.Get("addr")
	c.BindEnv("TOKEN")
	kubToken := c.Get("token")

	c.SetEnvPrefix("LOGSTASH")
	c.BindEnv("URI")
	logStashUri := c.Get("uri").(string)
	c.BindEnv("TAG")
	logStashTag := c.Get("tag").(string)
	c.BindEnv("NETWORK")
	logStashNetwork := c.Get("network").(string)
	c.BindEnv("LEVEL")
	logStasLevel := c.Get("level").(string)

	_, err := logger.New(logger.LoggingConfig{
		Tag:     logStashTag,
		ConnUri: logStashUri,
		Network: logStashNetwork,
		Level:   logStasLevel,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating logger: %v\n", err)
		os.Exit(2)
	}

	c.SetDefault("kubernetes_address", kubAddress)
	c.SetDefault("kubernetes_token", kubToken)
	c.SetDefault("dsn", fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable", pguser, pgpass, pgdatabase))
	c.SetDefault("cookie_secret", "qaBzlTixkx2c9S6i")
	c.SetDefault("http_addr", ":8888")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")

	c.AutomaticEnv()

	return c, nil
}

func main() {
	config, err := newConfig()
	if err != nil {
		logger.Fatal(err)
	}

	app, err := application.New(config)
	if err != nil {
		logger.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logger.Fatal(err)
	}

	serverAddress := config.Get("http_addr").(string)

	certFile := config.Get("http_cert_file").(string)
	keyFile := config.Get("http_key_file").(string)
	drainIntervalString := config.Get("http_drain_interval").(string)

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logger.Fatal(err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logger.Infoln("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logger.Fatal(err)
	}
}
