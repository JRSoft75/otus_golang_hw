package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/app"                          //nolint:depguard
	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/config"                       //nolint:depguard
	"github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/logger"                       //nolint:depguard
	internalhttp "github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/server/http"     //nolint:depguard
	storagePackage "github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage"       //nolint:depguard
	memorystorage "github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/JRSoft75/otus_golang_hw/hw12_13_14_15_16_calendar/internal/storage/sql"       //nolint:depguard
	"github.com/spf13/cobra"                                                                             //nolint:depguard
)

// var configFile string

// func init() {
//	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
//}

func main() {
	var configFile string
	var versionFlag bool

	rootCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Calendar service",
		Run: func(_ *cobra.Command, _ []string) {
			if versionFlag {
				printVersion()
				return
			}
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				log.Fatalf("failed to load config: %v", err)
			}

			logg, err := logger.New(cfg.Logger.Level)
			if err != nil {
				log.Fatalf("failed to initialize logger: %v", err)
			}
			defer logg.Sync()

			var storage storagePackage.Storage
			switch cfg.Storage.Type {
			case "in-memory":
				storage = memorystorage.NewInMemoryStorage()
			case "sql":
				db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=calendar sslmode=disable",
					cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password))
				if err != nil {
					log.Fatalf("failed to connect to database: %v", err)
				}
				storage = sqlstorage.NewSQLStorage(db)
			default:
				log.Fatalf("unsupported storage type: %s", cfg.Storage.Type)
			}

			logg.Info("calendar is running...")
			calendar := app.New(logg, storage)
			server := internalhttp.NewServer(
				logg,
				calendar,
				cfg.Server.Host,
				cfg.Server.Port,
				time.Duration(cfg.Server.ReadTimeout),
				time.Duration(cfg.Server.WriteTimeout),
			)

			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := server.Stop(ctx); err != nil {
					logg.Error("failed to stop http server: " + err.Error())
				}
			}()

			if err := server.Start(ctx); err != nil {
				logg.Error("failed to start http server: " + err.Error())
				cancel()
				os.Exit(1) //nolint:critic
			}
		},
	}
	// Флаг для указания пути к конфигурационному файлу
	rootCmd.Flags().StringVar(&configFile, "config", "", "path to config file (required)")

	// Флаг --version
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "print the version of the application")
	// err := rootCmd.MarkFlagRequired("config")
	// if err != nil {
	//	return
	// }

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("command execution failed: %v", err)
	}
}
