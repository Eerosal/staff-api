package mojang

import (
	"encoding/json"
	"fmt"
	"net/http"
	"staff-api/config"
	"staff-api/logger"
)

func FetchUsername(conf *config.Config, uuid string) (*string, error) {
	logger.Info(fmt.Sprintf("Fetching name for %v", uuid))

	resp, err := http.Get(fmt.Sprintf(conf.NameDataUrl, uuid))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make request: %v", resp.Status)
	}

	var profile ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &profile.Name, nil
}

// ProfileResponse We only need the name, so we can ignore the rest
type ProfileResponse struct {
	Name string `json:"name"`
}
