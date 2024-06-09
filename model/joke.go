package model

import (
	"encoding/json"
	"net/http"
)

// Joke struct to unmarshal JSON response
type Joke struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

// GetJoke fetches a random joke from the API
func GetJoke() (Joke, error) {
	var joke Joke
	resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {

		return joke, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		return joke, err
	}

	return joke, nil
}
