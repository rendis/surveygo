package main

import (
	"github.com/rendis/surveygo"
	"log"
	"os"
)

var nw = `{
      "order": 0,
      "nameId": "new_generals",
      "type": "checkbox",
      "label": "select one or more options",
      "value": {
        "options": [
          {
            "id": "1",
            "label": "option 1",
            "order": 0
          },
          {
            "id": "2",
            "label": "option 2",
            "order": 1
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

	err = s.Check(map[string][]any{"sub_generals": {"2"}, "generals": {"2"}, "generals-radio": {"1"}})
	if err != nil {
		log.Printf("%s", err)
	}

	// Add new question
	err = s.AddQuestionJson(nw)
	if err != nil {
		panic(err)
	}
}
