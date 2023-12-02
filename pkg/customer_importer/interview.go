// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bebsimax/teamwork_task/pkg/models"
)

type CustomerImporter struct {
	Log *log.Logger
}

func (ci *CustomerImporter) Load(path string) ([]*DomainCount, error) {
	f, err := os.Open(path)
	if err != nil {
		ci.Log.Printf("ERROR: open file %s: %s", path, err)

		return nil, fmt.Errorf("open file: %w", err)
	}

	defer f.Close()

	return ci.process(f, false)
}

func (ci *CustomerImporter) process(r io.Reader, failFast bool) ([]*DomainCount, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = models.CustomerFieldsPerLine
	curLine := 1

	m := map[string]int{}

	for {
		curLine++
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			ci.Log.Printf("ERROR line %d: read line: %s", curLine, err)
			if failFast {
				return nil, fmt.Errorf("read line %d: %w", curLine, err)
			}

			continue
		}

		c, err := models.NewCustomerFromLine(line)
		if err != nil {
			return nil, fmt.Errorf("new customer: %w", err)
		}

		err = c.Validate()
		if err != nil {
			ci.Log.Println(fmt.Sprintf("ERROR line %d: verify customer: %s", curLine, err))
			if failFast {
				return nil, fmt.Errorf("validate customer: %w", err)
			}

			continue
		}

		domain := strings.Split(c.Email, "@")[1]

		_, ok := m[domain]
		if !ok {
			m[domain] = 0
		}

		m[domain]++
	}

	return mapToSortedSlice(m), nil
}

type DomainCount struct {
	Domain string
	Count  int
}

func mapToSortedSlice(m map[string]int) []*DomainCount {
	s := make([]*DomainCount, len(m))

	i := 0

	for domain, count := range m {
		s[i] = &DomainCount{Domain: domain, Count: count}
		i++
	}

	if len(s) > 1 {
		sort.Slice(s, func(i, j int) bool { return s[i].Count > s[j].Count })
	}

	return s
}
