package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/GerardRodes/muzz-backend/internal/config"
	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/GerardRodes/muzz-backend/internal/httpserver"
	"github.com/GerardRodes/muzz-backend/internal/mariadb"
	"github.com/GerardRodes/muzz-backend/internal/session"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run handles all the initialization and defer callbacks
// moving it out of main ensures that defer calls desconnecting
// from sources will execute
func run() error {
	cfg := config.New()

	db, err := initDBHandle(cfg)
	if err != nil {
		return fmt.Errorf("cannot init db handler: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		} else {
			log.Println("closed db connection")
		}
	}()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.KVAddr})
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("cannot reach redis: %w", err)
		}
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Println(err)
		} else {
			log.Println("closed redis connection")
		}
	}()

	ss := session.NewSessionStorage(session.Config{
		RedisClient: rdb,
		Expiration:  time.Hour * 24 * 7,
	})

	if err := httpserver.Init(httpserver.Config{
		HTTPPort:        cfg.HTTPPort,
		Service:         domain.NewService(mariadb.NewRepo(db), ss),
		HandlersTimeout: time.Second * 5,
	}); err != nil {
		return fmt.Errorf("cannot init http server: %w", err)
	}

	return nil
}

func initDBHandle(cfg config.Config) (*sql.DB, error) {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.DBUser
	mysqlCfg.Passwd = cfg.DBPassword
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = cfg.DBAddr
	mysqlCfg.DBName = cfg.DBName

	connector, err := mysql.NewConnector(mysqlCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create connector: %w", err)
	}

	db := sql.OpenDB(connector)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot reach the database: %w", err)
	}

	return db, nil
}
