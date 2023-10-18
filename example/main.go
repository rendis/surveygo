package main

import (
	"fmt"
	"github.com/rendis/surveygo/v2"
	"log"
	"os"
)

var nw = `{
	"nameId": "email",
	"visible": true,
	"type": "email",
	"label": "Email",
	"required": true,
	"value": {
		"placeholder": "Type your email here",
		"allowedDomains": ["gmail.com", "yahoo.com", "hotmail.com"]
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
		"service_quality":      {"good"},
		"recommend_to_friends": {"definitely"},
		"additional_comments":  {"satisfied"},
	}

	// review answers
	resume, err := s.ReviewAnswers(ans)
	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("Resume: %+v\n", resume)

	// add new question
	err = s.AddQuestionJson(nw)
	if err != nil {
		log.Fatalf("Error adding question: %v", err)
	}

	// add new answer
	ans["email"] = []any{"test@yopmail.com"}

	// review answers
	resume, err = s.ReviewAnswers(ans)

	if err != nil {
		log.Fatalf("Error checking survey: %v", err)
	}

	fmt.Printf("Resume: %+v\n", resume)
}
