package main

import (
	"fmt"
	"log"
	"os"

	customerimporter "github.com/bebsimax/teamwork_task/pkg/customer_importer"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	c := customerimporter.CustomerImporter{Log: logger}
	s, err := c.Load("customers.csv")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		fmt.Println(*s[i])
	}
}
