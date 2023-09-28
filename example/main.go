package main

import (
	"fmt"
	"github.com/rendis/surveygo/v2"
	"log"
	"os"
)

var nw = `{
	"nameId": "new_generals",
	"visible": true,
	"type": "checkbox",
	"label": "select one or more options",
	"groupNameId": "generals-info",
	"value": {
	  "options": [
		{
		  "nameId": "qwe1",
		  "label": "option 1"
		},
		{
		  "nameId": "qwe2",
		  "label": "option 2"
		}
	  ]
	}
  }`

func main() {
	f, err := os.ReadFile("example/s1.json")
	if err != nil {
		panic(err)
	}

	s, err := surveygo.ParseBytes(f)
	if err != nil {
		panic(err)
	}

	ans := map[string][]any{
		"generals":       {"qwq2"},
		"generals-radio": {"qwr1"},
		"hidden-q":       {"qwh2"},
	}

	resume, err := s.Check(ans)
	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("Resume: %+v\n", resume)

	// Add new question
	err = s.AddQuestionJson(nw)
	if err != nil {
		log.Fatalf("Error adding question: %v", err)
	}
}
