package models

import (
	"fmt"
	"regexp"
)

const CustomerFieldsPerLine = 5

type Customer struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

func (c *Customer) Validate() error {
	expression := `\b[\w\.-]+@[\w\.-]+\.\w{2,4}\b`
	ok, err := regexp.MatchString(expression, c.Email)
	if err != nil {
		return fmt.Errorf("process regexp %s", expression)
	}

	if !ok {
		return fmt.Errorf("email: %s does not match regexp", c.Email)
	}

	return nil
}

func NewCustomerFromLine(line []string) (*Customer, error) {
	if len(line) != CustomerFieldsPerLine {
		return nil, fmt.Errorf("unexpected number of fields in line, expected: %d, got: %d", CustomerFieldsPerLine, len(line))
	}

	return &Customer{FirstName: line[0],
		LastName:  line[1],
		Email:     line[2],
		Gender:    line[3],
		IPAddress: line[4]}, nil
}
