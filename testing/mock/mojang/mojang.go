package main

import (
	"fmt"
	"net/http"
)

func main() {
	// The real API uses path parameters, but we're using query parameters here since
	// they are easier to handle with the default http library
	http.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		uuid := r.URL.Query().Get("uuid")
		if uuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := ""
		if uuid == "40a1b924-e5a8-4444-9fec-73db15ee7c8d" {
			name = "player1"
		} else if uuid == "c806a034-b626-4336-a1e2-37b902888bf5" {
			name = "player2"
		} else if uuid == "d0b765a7-87d6-49fb-96d4-17e1cdb6ca2e" {
			name = "player3"
		} else if uuid == "c706ec9f-1d29-4cf6-a849-ba5830b5cb41" {
			name = "player4"
		} else if uuid == "b1a2291f-3955-431e-b274-2388f85d3b63" {
			name = "player5"
		} else if uuid == "ef7eb665-a3ac-40a6-b9a4-1100f60b28cd" {
			name = "player6"
		}

		if len(name) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, err := w.Write([]byte(`{"name":"` + name + `"}`))
		if err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe("staff-api-mojang-mock:8080", nil)
	if err != nil {
		panic(err)
	}
}
