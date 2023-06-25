package fetch

import (
	"encoding/json"
	"fmt"
	"sort"
	"staff-api/common"
	"staff-api/config"
	"staff-api/fetch/avatar"
	"staff-api/fetch/luckperms"
	"staff-api/fetch/mojang"
	"staff-api/logger"
	"sync"
)

type userEntry struct {
	Uuid   string   `json:"uuid"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	Avatar *string  `json:"avatar,omitempty"`
}

type userResponse struct {
	Users []userEntry `json:"users"`
}

type Result struct {
	Result            []byte
	ResultWithAvatars []byte
}

func Fetch(conf *config.Config) (*Result, error) {
	groups, err := luckperms.FetchUserDisplayGroups(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group meta: %w", err)
	}

	uuids := common.GetMapKeys(groups)
	var mu sync.Mutex

	names := make(map[string]string)
	common.RunRateLimited(uuids, 3, func(uuid string) {
		name, err := mojang.FetchUsername(conf, uuid)
		if err != nil {
			logger.Warn(fmt.Sprintf("failed to fetch name for %v: %v", uuid, err))
			unknown := "Unknown"
			name = &unknown
		}

		if name != nil {
			mu.Lock()
			names[uuid] = *name
			mu.Unlock()
		}
	})

	avatarUrls := make(map[string]string)
	common.RunRateLimited(uuids, 4, func(uuid string) {
		avatarUrl, err := avatar.FetchAvatar(conf, uuid)
		if err != nil {
			logger.Warn(fmt.Sprintf("failed to fetch avatar for %v: %v", uuid, err))
			return
		}

		if avatarUrl != nil {
			mu.Lock()
			avatarUrls[uuid] = *avatarUrl
			mu.Unlock()
		}
	})

	users := make([]userEntry, 0, len(groups))

	for uuid, groupEntry := range groups {
		name, ok := names[uuid]
		if !ok {
			logger.Warn(fmt.Sprintf("failed to find name for %v", uuid))
			continue
		}

		entry := userEntry{
			Uuid:   uuid,
			Name:   name,
			Groups: []string{groupEntry.Group},
		}

		avatarUrl, ok := avatarUrls[uuid]
		if ok {
			entry.Avatar = &avatarUrl
		}

		sort.Strings(entry.Groups)

		users = append(users, entry)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Uuid < users[j].Uuid
	})

	response := userResponse{
		Users: users,
	}

	withAvatars, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	for i := 0; i < len(users); i++ {
		users[i].Avatar = nil
	}

	withoutAvatars, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return &Result{
		Result:            withoutAvatars,
		ResultWithAvatars: withAvatars,
	}, nil
}
