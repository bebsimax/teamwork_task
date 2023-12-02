package customerimporter

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	for _, tt := range []struct {
		name     string
		data     []string
		expected []*DomainCount
		fails    bool
		errMsg   string
	}{
		{
			name:     "no records",
			expected: []*DomainCount{},
		},
		{
			name: "single record",
			data: []string{
				"fname,lname,mail@sample.com,g,ip",
			},
			expected: []*DomainCount{
				{Domain: "sample.com", Count: 1},
			},
		},
		{
			name: "multiple records",
			data: []string{
				"fname,lname,mail@sample.com,g,ip",
				"fname,lname,mail@sample2.com,g,ip",
				"fname,lname,mail@sample3.com,g,ip",
			},
			expected: []*DomainCount{
				{Domain: "sample.com", Count: 1},
				{Domain: "sample2.com", Count: 1},
				{Domain: "sample3.com", Count: 1},
			},
		},
		{
			name: "multiple same domains",
			data: []string{
				"fname,lname,mail@sample.com,g,ip",
				"fname,lname,mail@sample.com,g,ip",
				"fname,lname,mail@sample.com,g,ip",
			},
			expected: []*DomainCount{
				{Domain: "sample.com", Count: 3},
			},
		},
		{
			name: "multiple domains",
			data: []string{
				"fname,lname,mail@sample.com,g,ip",
				"fname,lname,mail@sample1.com,g,ip",
				"fname,lname,mail@sample1.com,g,ip",
				"fname,lname,mail@sample2.com,g,ip",
				"fname,lname,mail@sample2.com,g,ip",
			},
			expected: []*DomainCount{
				{Domain: "sample1.com", Count: 2},
				{Domain: "sample2.com", Count: 2},
				{Domain: "sample.com", Count: 1},
			},
		},
		{
			name: "error: unexpected number of fields",
			data: []string{
				"fname,lname",
			},
			fails:  true,
			errMsg: "wrong number of fields",
		},
		{
			name: "error: unparsable email",
			data: []string{
				"fname,lname,mailsample.com,g,ip",
			},
			fails:  true,
			errMsg: "does not match regexp",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)

			mock := &MockReader{data: tt.data}
			ci := CustomerImporter{Log: logger}

			c, err := ci.process(mock, true)
			if tt.fails {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errMsg)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, c)
		})
	}
}

func TestMapToSortedSlice(t *testing.T) {
	for _, tt := range []struct {
		name          string
		m             map[string]int
		expectedSlice []*DomainCount
	}{
		{
			name:          "empty",
			expectedSlice: []*DomainCount{},
		},
		{
			name: "single entry",
			m: map[string]int{
				"d1": 1,
			},
			expectedSlice: []*DomainCount{{Domain: "d1", Count: 1}},
		},
		{
			name: "multiple entries",
			m: map[string]int{
				"d1": 1,
				"d2": 2,
				"d3": 3,
			},
			expectedSlice: []*DomainCount{
				{Domain: "d3", Count: 3},
				{Domain: "d2", Count: 2},
				{Domain: "d1", Count: 1},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := mapToSortedSlice(tt.m)
			assert.Equal(t, tt.expectedSlice, s)
		})
	}
}

type MockReader struct {
	data  []string
	index int
}

func (f *MockReader) Read(p []byte) (int, error) {
	if f.index == len(f.data) {
		return 0, io.EOF
	}

	b := []byte(f.data[f.index] + "\n")
	copy(p[:], b)
	f.index++

	return len(b), nil
}
