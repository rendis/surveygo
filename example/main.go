package main

import (
	"encoding/json"
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

	s, err := surveygo.ParseFromBytes(f)
	if err != nil {
		panic(err)
	}

	// -------------  answer individual questions -------------- //

	ans1 := map[string][]any{
		"event_rating":       {"good"},
		"favorite_game":      {"zombie_apocalypse"},
		"would_attend_again": {"would_attend_again_yes"},
		"name":               {"John Doe"},
		"email":              {"john@example.com"},
		"phone_codes":        {"+65"},
		"phone-number":       {"12345678"},
		"birth_date":         {"2021-01-01"},
	}

	// review answers
	resume1, err := s.ReviewAnswers(ans1)
	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("\nResume: %+v\n", resume1)

	// add new question
	err = s.AddQuestionJson(newQuestion)
	if err != nil {
		log.Fatalf("Error adding question: %v", err)
	}

	// add new question to group
	if err = s.AddQuestionToGroup("favourite_game_song", "grp-general", -1); err != nil {
		log.Fatalf("Error adding question to group: %v", err)
	}

	// add new answer
	ans1["favourite_game_song"] = []any{"song_cyberworld_2078"}

	// review answers
	resume2, err := s.ReviewAnswers(ans1)

	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("\nResume: %+v\n", resume2)

	// -------------  answer grouped questions -------------- //

	var ans2 = map[string][]map[string][]any{
		"grp-general": {
			{
				"event_rating":        {"good"},
				"favorite_game":       {"zombie_apocalypse"},
				"would_attend_again":  {"would_attend_again_yes"},
				"favourite_game_song": {"song_cyberworld_2078"},
			},
			{
				"event_rating":        {"good"},
				"favorite_game":       {"zombie_apocalypse"},
				"would_attend_again":  {"would_attend_again_yes"},
				"favourite_game_song": {"song_cyberworld_2078"},
			},
		},
		"grp-feedback": {
			{
				"name":  {"John Doe"},
				"email": {"john@example.com"},
			},
		},
	}

	var ans2Casted = make(map[string][]any)
	for k, v := range ans2 {
		ans2Casted[k] = []any{v}
	}

	// add new answer
	ans2Casted["event_rating_improvement_opinion"] = []any{
		"the event was great, but the food was not good enough",
	}

	// review answers
	resume3, err := s.ReviewAnswers(ans2Casted)

	if resume3.InvalidAnswers != nil {
		log.Print("Error checking survey for grouped answers")
		for _, v := range resume3.InvalidAnswers {
			fmt.Printf(" - Invalid answer: %s\n", v)
		}
	}

	// resume to json
	var surveyResumeJson []byte
	surveyResumeJson, err = json.Marshal(resume3)
	if err != nil {
		log.Fatalf("Error marshaling survey resume: %v", err)
	}

	fmt.Printf("\nSurvey Resume: %s\n", surveyResumeJson)
}
