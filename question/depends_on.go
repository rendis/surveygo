package question

type DependsOn struct {
	QuestionNameId string `json:"questionNameId" bson:"questionNameId" validate:"required,validNameId"`
	OptionNameId   string `json:"optionNameId" bson:"optionNameId" validate:"required,validNameId"`
}
