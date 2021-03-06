package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/soulmonk/go-grpc-http-rest-microservice-tutorial/pkg/logger"

	//"github.com/jmoiron/sqlx" TODO
	"log"

	// postgres driver
	_ "github.com/lib/pq"

	"github.com/soulmonk/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"github.com/soulmonk/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"
	"github.com/soulmonk/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
)

type PG struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	Db PG

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	// new line
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.Db.Host, "db-host", "", "Database host")
	flag.StringVar(&cfg.Db.Port, "db-port", "", "Database port")
	flag.StringVar(&cfg.Db.User, "db-user", "", "Database user")
	flag.StringVar(&cfg.Db.Password, "db-password", "", "Database password")
	flag.StringVar(&cfg.Db.Dbname, "db-name", "", "Database name")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP getaway: '%s'", cfg.HTTPPort)
	}

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Db.Host, cfg.Db.Port,
		cfg.Db.User, cfg.Db.Password, cfg.Db.Dbname)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	v1API := v1.NewToDoServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
