package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type configLuckperms struct {
	ConnectionString string
}

type Config struct {
	Groups              map[string]struct{}
	Interval            time.Duration
	ImageUrlFn          func(uuid string) *string
	ImageUpdateInterval time.Duration
	Luckperms           configLuckperms
	BindAddress         string
	NameDataUrl         string
}

func ParseConfig() (*Config, error) {
	conf := Config{
		Luckperms: configLuckperms{},
	}
	conf.Groups = make(map[string]struct{})
	for _, group := range strings.Split(os.Getenv("STAFF_API_GROUPS"), ",") {
		conf.Groups[group] = struct{}{}
	}

	imageUrl := os.Getenv("STAFF_API_IMAGE_URL")
	conf.ImageUrlFn = func(uuid string) *string {
		if len(imageUrl) == 0 {
			return nil
		}
		url := fmt.Sprintf(imageUrl, string(uuid))
		return &url
	}
	imageUpdateIntervalSec := os.Getenv("STAFF_API_IMAGE_UPDATE_INTERVAL")
	if len(imageUpdateIntervalSec) == 0 {
		imageUpdateIntervalSec = "60"
	}
	imageUpdateInterval, err := time.ParseDuration(imageUpdateIntervalSec + "s")
	if err != nil {
		return nil, err
	}
	if imageUpdateInterval < 1*time.Minute {
		return nil, fmt.Errorf("image update interval must be at least 1 minute")
	}

	updateIntervalSec := os.Getenv("STAFF_API_UPDATE_INTERVAL")
	if len(updateIntervalSec) == 0 {
		updateIntervalSec = "60"
	}
	updateInterval, err := time.ParseDuration(updateIntervalSec + "s")
	if err != nil {
		return nil, err
	}
	if updateInterval < 1*time.Minute {
		return nil, fmt.Errorf("update interval must be at least 1 minute")
	}
	conf.Interval = updateInterval

	conf.Luckperms.ConnectionString = os.Getenv("STAFF_API_LUCKPERMS_CONNECTION_STRING")

	conf.BindAddress = os.Getenv("STAFF_API_REST_HTTP_ADDRESS") + ":" + os.Getenv("STAFF_API_REST_HTTP_PORT")

	conf.NameDataUrl = os.Getenv("STAFF_API_NAME_DATA_URL")
	if len(conf.NameDataUrl) == 0 {
		conf.NameDataUrl = "https://sessionserver.mojang.com/session/minecraft/profile/%s"
	}

	return &conf, nil
}
