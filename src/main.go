package main

import (
	"bytes"
	"fmt"
	"net/http"
	"staff-api/config"
	"staff-api/fetch"
	"staff-api/fetch/luckperms"
	"staff-api/logger"
	"sync"
	"time"
)

func main() {
	var rwm sync.RWMutex

	var resultUpdateTime time.Time
	var result *fetch.Result

	conf, err := config.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	err = luckperms.WaitForDBConnection(conf, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	result, err = fetch.Fetch(conf)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to fetch: %v", err))
	}

	go func() {
		t := time.NewTicker(conf.Interval)
		for {
			<-t.C
			newResult, err := fetch.Fetch(conf)
			if err != nil {
				logger.Warn(fmt.Sprintf("Failed to fetch: %v", err))
				continue
			}

			rwm.Lock()

			now := time.Now()
			if result == nil || !bytes.Equal(newResult.ResultWithAvatars, result.ResultWithAvatars) {
				result = newResult
				logger.Info("Result updated")
				resultUpdateTime = now
			} else {
				logger.Info("Result unchanged")
			}

			rwm.Unlock()
		}
	}()

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, no-cache")

		var response []byte = nil

		ifModifiedStr := r.Header.Get("If-Modified-Since")
		rwm.RLock()
		if result == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("Internal server error"))
			if err != nil {
				logger.Warn(fmt.Sprintf("Failed to write response: %v", err))
				return
			}
			rwm.RUnlock()
			return
		}
		w.Header().Set("Content-Type", "application/json")

		if ifModified, err := time.Parse(
			http.TimeFormat,
			ifModifiedStr,
		); err != nil || resultUpdateTime.After(ifModified) {
			avatarsParam := r.URL.Query().Get("avatars")
			if avatarsParam == "true" {
				response = make([]byte, len(result.ResultWithAvatars))
				copy(response, result.ResultWithAvatars)
			} else {
				response = make([]byte, len(result.Result))
				copy(response, result.Result)
			}
		}
		rwm.RUnlock()

		if response == nil {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Last-Modified", resultUpdateTime.Format(http.TimeFormat))
		_, err := w.Write(response)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to write response: %v", err))
			return
		}
	})

	logger.Info(fmt.Sprintf("Listening on %v", conf.BindAddress))

	if err := http.ListenAndServe(conf.BindAddress, nil); err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}
}
