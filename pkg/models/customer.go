package models

import (
	"fmt"
	"regexp"
)

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

func NewCustomer(row []string) *Customer {
	return &Customer{FirstName: row[0],
		LastName:  row[1],
		Email:     row[2],
		Gender:    row[3],
		IPAddress: row[4]}
}
