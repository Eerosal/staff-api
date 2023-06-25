package avatar

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"staff-api/config"
	"staff-api/logger"
	"time"
)

func FetchAvatar(conf *config.Config, uuid string) (*string, error) {
	url := conf.ImageUrlFn(uuid)
	if url == nil {
		return nil, nil
	}

	logger.Info(fmt.Sprintf("Fetching avatar for %v", uuid))

	avatarCacheMutex.RLock()
	cacheEntry, ok := avatarCache[uuid]
	var cachedAvatar *string = nil
	var cachedAvatarModifiedTime time.Time
	if ok && time.Since(cacheEntry.LastFetchTime) < conf.ImageUpdateInterval {
		cachedAvatar = cacheEntry.clone().LastAvatar
		cachedAvatarModifiedTime = cacheEntry.LastFetchTime
	}
	avatarCacheMutex.RUnlock()

	if cachedAvatar != nil {
		logger.Info(fmt.Sprintf("Avatar for %v not modified (cached)", uuid))
		return cacheEntry.LastAvatar, nil
	}

	request, err := http.NewRequest(http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if cachedAvatar != nil {
		request.Header.Set("If-Modified-Since", cachedAvatarModifiedTime.Format(time.RFC1123))
	}

	logger.Info(fmt.Sprintf("Fetching avatar for %v from %v", uuid, *url))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusNotModified {
		logger.Info(fmt.Sprintf("Avatar for %v not modified (304)", uuid))
		return cacheEntry.LastAvatar, nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make request: %v", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if !isValidPng(data) {
		return nil, fmt.Errorf("invalid PNG data")
	}

	avatar := toBase64Url(data)
	avatarCacheMutex.Lock()
	avatarCache[uuid] = &avatarCacheEntry{
		LastFetchTime: time.Now(),
		LastAvatar:    avatar,
	}
	avatarCacheMutex.Unlock()

	return avatar, nil
}

func isValidPng(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	return string(data[:8]) == "\x89PNG\r\n\x1a\n"
}

func toBase64Url(data []byte) *string {
	if data == nil {
		return nil
	}

	str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	return &str
}
