package luckperms

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"staff-api/config"
	"time"
)

func getDBConnection(conf *config.Config) (*sql.DB, error) {
	return sql.Open("mysql", conf.Luckperms.ConnectionString)
}

func WaitForDBConnection(conf *config.Config, timeout time.Duration) error {
	start := time.Now()
	for {
		db, err := getDBConnection(conf)
		if err != nil {
			return err
		}

		err = db.Ping()
		if err != nil {
			if time.Since(start) > timeout {
				return err
			}
			time.Sleep(1 * time.Second)
			_ = db.Close()
			continue
		}

		break
	}

	return nil
}
