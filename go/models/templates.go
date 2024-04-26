package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "github.com/mjarkk/whatsapp-dev/go/db"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Name string `json:"name"`
	// FIXME different types of headers
	Header                *string                `json:"header"`
	Body                  string                 `json:"body"`
	Footer                *string                `json:"footer"`
	TemplateCustomButtons []TemplateCustomButton `json:"templateCustomButtons"`
}

type TemplateCustomButton struct {
	gorm.Model
	TemplateID uint   `json:"templateId"`
	Text       string `json:"text"`
}

var TemplateVriableRegex = regexp.MustCompile(`\{\{\s*\d+\s*\}\}`)

func (t *Template) CreateCustomButton(text string) error {
	button := TemplateCustomButton{
		TemplateID: t.ID,
		Text:       text,
	}
	err := DB.Create(&button).Error
	if err != nil {
		return err
	}

	t.TemplateCustomButtons = append(t.TemplateCustomButtons, button)
	return nil
}

// Variables returns all the variables in the template
func Variables(input string) []int {
	seenVariables := map[int]struct{}{}

	variables := TemplateVriableRegex.FindAllString(input, -1)
	for _, variable := range variables {
		variableNumber, err := strconv.Atoi(strings.Trim(variable, "{} "))
		if err != nil {
			continue
		}

		seenVariables[variableNumber] = struct{}{}
	}

	resp := []int{}
	for number := range seenVariables {
		resp = append(resp, number)
	}

	return resp
}

func ReplaceVariables(input string, varValues []string) string {
	return TemplateVriableRegex.ReplaceAllStringFunc(input, func(variablePlaceholder string) string {
		variableNumber, err := strconv.Atoi(strings.Trim(variablePlaceholder, "{} "))
		if err != nil {
			return variablePlaceholder
		}
		if variableNumber <= 0 {
			return variablePlaceholder
		}
		if variableNumber > len(varValues) {
			return variablePlaceholder
		}
		return varValues[variableNumber-1]
	})
}

func validateVariables(input string) error {
	seenNumbers := map[int]struct{}{}

	variables := TemplateVriableRegex.FindAllString(input, -1)
	for _, variable := range variables {
		variableNumber, err := strconv.Atoi(strings.Trim(variable, "{} "))
		if err != nil {
			return fmt.Errorf("variable %s is not a number", variable)
		}
		if variableNumber == 0 {
			return fmt.Errorf("variable %s, this is not allowed to be zero", variable)
		}

		seenNumbers[variableNumber] = struct{}{}
	}

	for key := range seenNumbers {
		if key == 1 {
			continue
		}
		_, ok := seenNumbers[key-1]
		if !ok {
			return fmt.Errorf("variable {{ %d }} is missing", key-1)
		}
	}

	return nil
}

func (t *Template) Validate() error {
	if t.Header != nil && *t.Header == "" {
		t.Header = nil
	}
	if t.Footer != nil && *t.Footer == "" {
		t.Footer = nil
	}

	err := validateVariables(t.Body)
	if err != nil {
		return fmt.Errorf("body: %s", err.Error())
	}

	if t.Header != nil {
		err = validateVariables(*t.Header)
		if err != nil {
			return fmt.Errorf("header: %s", err.Error())
		}
	}

	return nil
}
