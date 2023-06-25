package luckperms

import (
	"database/sql"
	"fmt"
	"staff-api/config"
	"staff-api/logger"
)

func FetchUserDisplayGroups(conf *config.Config) (map[string]UserGroupEntry, error) {
	logger.Info("Fetching user groups")
	db, err := getDBConnection(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, err := db.Query("SELECT uuid, permission FROM luckperms_user_permissions WHERE permission LIKE 'group.%' AND value <> 0")
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	_, includeAll := conf.Groups["*"]

	results := make(map[string]UserGroupEntry, 0)
	for rows.Next() {
		var uuid string
		var permission string
		err = rows.Scan(&uuid, &permission)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		group := permission[6:]

		if !includeAll {
			if _, ok := conf.Groups[group]; !ok {
				continue
			}
		}

		results[uuid] = UserGroupEntry{
			UniqueId: uuid,
			Group:    group,
		}
	}

	return results, nil
}

type UserGroupEntry struct {
	UniqueId string
	Group    string
}
