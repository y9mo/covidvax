package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	goflag "flag"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/y9mo/covidvax"
	"github.com/y9mo/covidvax/api"
	"github.com/y9mo/covidvax/repository"
)

/// Prefix for environments variables
const prefix string = "COVIDVAX"

type Config struct {
	Development  bool   `mapstructure:"dev"`
	PgConnection string `mapstructure:"pg-connection"`
	Listen       string `mapstructure:"listen"`
}

func GetConfig() (Config, error) {
	pflag.Bool("dev", false, "enable development mode")
	pflag.String("pg-connection",
		"host=127.0.0.1 port=5432 user=admin dbname=covidvax password=admin-pwd sslmode=disable",
		"postgresql connection string")
	pflag.String("listen", ":8080", "listen address")

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()

	var config Config
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return config, fmt.Errorf("unable to parse flags, %w", err)
	}
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("flags unmarshalling failed: %w", err)
	}

	return config, nil
}

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Printf("unable to load data at startup: %s\n", err)
		os.Exit(1)
	}

	logger, err := initLog(config.Development)
	if err != nil {
		fmt.Printf("unable to init logger: %s\n", err)
		os.Exit(1)
	}
	logger.Sugar().Debugf("%+v", config)

	logger.Sugar().Debugf("connecting to %s", config.PgConnection)
	db, err := gorm.Open("postgres", config.PgConnection)
	if err != nil {
		logger.Sugar().Fatalf("failed to connect to the database: %s", err)
	}

	pr := repository.NewPatients(db, logger)
	tcr := repository.NewTreatmentCenters(db, logger)
	ar := repository.NewAppointments(db, logger)
	abr := repository.NewAppointmentBookings(db, logger)

	router, err := api.Setup(logger, pr, tcr, ar, abr)
	if err != nil {
		logger.Sugar().Fatalf("router setup: %s", err)
	}

	srv := http.Server{Addr: config.Listen, Handler: router}

	logger.Sugar().Infof("listening on %s", config.Listen)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar().Fatalf("listen: %s", err)
	}
}

func initLog(development bool) (logger *zap.Logger, err error) {
	if development {
		logger, err = zap.NewDevelopmentConfig().Build()
	} else {
		logger, err = zap.NewProductionConfig().Build()
	}
	if err != nil {
		return nil, err
	}
	logger = logger.With(zap.String("version", covidvax.Version))
	return logger, nil
}
