package main

import (
	"fmt"
	"github.com/rendis/surveygo/v2"
	"log"
	"os"
)

var newQuestion = `{
        "nameId": "favourite_game_song",
        "visible": true,
        "type": "single_select",
        "label": "What is your favourite game song?",
        "required": false,
        "value": {
            "options": [
            {
                "nameId": "song_cyberworld_2078",
                "label": "Cyberworld 2078"
            },
            {
                "nameId": "song_zombie_apocalypse",
                "label": "Zombie Apocalypse"
            },
            {
                "nameId": "song_fortress_siege",
                "label": "Fortress Siege"
            },
            {
                "nameId": "song_rocket_adventure",
                "label": "Rocket Adventure"
            }
            ]
        }
    }`

func main() {
	f, err := os.ReadFile("example/survey.json")
	if err != nil {
		panic(err)
	}

	s, err := surveygo.ParseBytes(f)
	if err != nil {
		panic(err)
	}

	ans := map[string][]any{
		"event_rating":       {"good"},
		"favorite_game":      {"zombie_apocalypse"},
		"would_attend_again": {"would_attend_again_yes"},
		"name":               {"John Doe"},
		"email":              {"john@example.com"},
		"phone_codes":        {"+65"},
		"phone-number":       {"12345678"},
	}

	// review answers
	resume, err := s.ReviewAnswers(ans)
	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("\nResume: %+v\n", resume)

	// add new question
	err = s.AddQuestionJson(newQuestion)
	if err != nil {
		log.Fatalf("Error adding question: %v", err)
	}

	// add new answer
	ans["favourite_game_song"] = []any{"song_cyberworld_2078"}

	// review answers
	resume, err = s.ReviewAnswers(ans)

	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("\nResume: %+v\n", resume)
}
