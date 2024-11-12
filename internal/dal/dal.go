package dal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"protchain/internal/config"
	"sync"
	"time"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DAL struct {
	SqlDB *bun.DB
}

// connectSQlDAL ...
func connectSQLDAL(config *config.Config) *bun.DB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DatabaseUrl)))
	sqlDB.SetMaxOpenConns(config.MaximumDBConn)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(time.Second * 2)
	Conn := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if os.Getenv("ENVIRONMENT") != "live" && os.Getenv("ENVIRONMENT") != "production" {
		Conn.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return Conn
}

func CreateTables(Conn *bun.DB) error {
	models := []interface{}{}

	var wg sync.WaitGroup
	for _, model := range models {

		wg.Add(1)

		go func(m any) {
			defer wg.Done()
			_, err := Conn.NewCreateTable().
				IfNotExists().
				Model(m).Exec(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
		}(model)

	}

	wg.Wait()
	return nil
}

type TableIndex struct {
	IndexName string
	Columns   []string
}

func CreateIndex(Conn *bun.DB) error {
	index := map[any]TableIndex{}

	var wg sync.WaitGroup
	for model, mIndex := range index {

		wg.Add(1)

		go func(m any, index TableIndex) {
			defer wg.Done()
			_, err := Conn.NewCreateIndex().
				IfNotExists().
				Index(index.IndexName).
				Column(index.Columns...).
				Model(m).Exec(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
		}(model, mIndex)

	}

	wg.Wait()
	return nil
}

func NewDAL(cfg *config.Config) *DAL {
	sqlDB := connectSQLDAL(cfg)
	if err := CreateTables(sqlDB); err != nil {
		log.Fatal(err)
	}

	return &DAL{
		SqlDB: sqlDB,
	}
}
