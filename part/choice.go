package part

import (
	"encoding/json"
	"fmt"
)

type Choice struct {
	Options []Option `json:"options"`
}

type Option struct {
	ID           *string    `json:"id"`           // Required
	Label        *string    `json:"label"`        // Required
	Order        *int       `json:"order"`        // Optional
	SubQuestions []Question `json:"subQuestions"` // Optional
}

func (c *Choice) UnmarshalJSON(b []byte) error {
	type alias Choice
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*c = Choice(a)
	return nil
}

func ChoiceUnmarshallValidator(c *Choice) error {
	if c == nil {
		return fmt.Errorf("choice is nil")
	}
	if c.Options == nil || len(c.Options) == 0 {
		return fmt.Errorf("invalid choice format, options is not defined")
	}
	return nil
}
