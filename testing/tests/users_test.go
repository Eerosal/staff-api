package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestUsersEndpoint(t *testing.T) {
	err := waitForServer()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i += 1 {
		withAvatars := i == 1

		reqUrl := "http://staff-api:8884/api/users"
		if withAvatars {
			reqUrl += "?avatars=true"
		}

		resp, err := http.Get(reqUrl)
		if err != nil {
			t.Fatal(err)
		}
		// This is fine since there are only 2 iterations
		//goland:noinspection GoDeferInLoop
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		if resp.StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		// Expect response Content-Type to be application/json
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Fatalf("Expected Content-Type to be application/json, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		parsed := UsersEndpointResponseWithAvatars{
			Users: make([]UserEntryWithAvatars, 0),
		}
		// Parse as map[string]interface{} first to avoid errors due to the avatar field being null
		var parsedMap map[string]interface{}
		err = json.Unmarshal(body, &parsedMap)
		if err != nil {
			t.Fatal(err)
		}

		usersRaw, ok := parsedMap["users"]
		if !ok {
			t.Fatal("Expected users field in response")
		}

		usersArr, ok := usersRaw.([]interface{})
		if !ok {
			t.Fatal("Expected users field to be an array")
		}

		// Workaround for the avatar field not being present in the UserEntry struct
		for _, userEntry := range usersArr {
			userMap, ok := userEntry.(map[string]interface{})
			avatarRaw, ok := userMap["avatar"]
			var avatar *string = nil
			if ok && avatarRaw != nil {
				avatarStr, ok := avatarRaw.(string)
				if ok {
					avatar = &avatarStr
				}
			}

			// Hacky workaround: marshal the userMap to JSON and then unmarshal it into a UserEntry (without avatar)
			// Then convert to UserEntryWithAvatars and add the avatar field
			marshal, err := json.Marshal(userMap)
			if err != nil {
				t.Fatal(err)
			}

			var user UserEntry
			err = json.Unmarshal(marshal, &user)
			if err != nil {
				t.Fatal(err)
			}

			var userWithAvatar = UserEntryWithAvatars{
				Uuid:   user.Uuid,
				Name:   user.Name,
				Groups: user.Groups,
				Avatar: avatar,
			}

			parsed.Users = append(parsed.Users, userWithAvatar)
		}

		users := parsed.Users

		if len(users) != 4 {
			t.Fatalf("Expected 4 users, got %d", len(users))
		}

		if users[0].Uuid != "40a1b924-e5a8-4444-9fec-73db15ee7c8d" {
			t.Fatal("Expected user 0 to have uuid 40a1b924-e5a8-4444-9fec-73db15ee7c8d (lowest uuid in alphabetical order)")
		}
		if users[3].Uuid != "ef7eb665-a3ac-40a6-b9a4-1100f60b28cd" {
			t.Fatal("Expected user 3 to have uuid ef7eb665-a3ac-40a6-b9a4-1100f60b28cd (highest uuid in alphabetical order)")
		}

		for _, expectedPlayer := range []string{"player1", "PLAYER2", "player3", "player6"} {
			found := false
			for _, user := range users {
				if user.Name == expectedPlayer {
					found = true
				}
			}
			if !found {
				t.Fatalf("Expected user %v to be present", expectedPlayer)
			}
		}

		for _, unexpectedPlayer := range []string{"player4", "player5"} {
			found := false
			for _, user := range users {
				if user.Name == unexpectedPlayer {
					found = true
				}
			}
			if found {
				t.Fatalf("Expected user %v to not be present", unexpectedPlayer)
			}
		}

		for _, user := range users {
			if user.Name == "" {
				t.Fatal("Expected user to have a name")
			}
			if len(user.Groups) != 1 {
				t.Fatal("Expected user to have exactly 1 group")
			}

			group := user.Groups[0]
			if group != "test1" && group != "test3" {
				t.Fatalf("Expected user to have test1 or test3, got %v", group)
			}

			if user.Name == "player1" && group != "test1" {
				t.Fatalf("Expected user %v to be in test1", user.Name)
			}

			if user.Name == "player6" && group != "test3" {
				t.Fatalf("Expected user %v to be in test3", user.Name)
			}

			if user.Name == "player4" || user.Name == "player5" {
				t.Fatalf("Expected user %v not to be displayed", user.Name)
			}

			if withAvatars {
				// Only player1 should have an avatar
				shouldHaveAvatar := user.Name == "player1"
				if shouldHaveAvatar && user.Avatar == nil {
					t.Fatalf("Expected user %v to have an avatar", user.Name)
				} else if !shouldHaveAvatar && user.Avatar != nil {
					t.Fatalf("Expected user %v to not have an avatar", user.Name)
				}
			} else {
				if user.Avatar != nil {
					t.Fatalf("Expected user %v to not have an avatar", user.Name)
				}
			}
		}
	}
}

type UserEntry struct {
	Uuid   string   `json:"uuid"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

// Apparently Unmarshal doesn't support omitempty, while Marshal does. Nice.
type UserEntryWithAvatars struct {
	Uuid   string   `json:"uuid"`
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	Avatar *string  `json:"avatar,omitempty"`
}

type UsersEndpointResponseWithAvatars struct {
	Users []UserEntryWithAvatars `json:"users"`
}

func waitForServer() error {
	for i := 0; i < 10; i += 1 {
		resp, err := http.Get("http://staff-api:8884/api/users")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("server did not start in time")
}
