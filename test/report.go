package test

import "log"

type Report struct {
	Total int
}

func (report *Report) Run(name string, testCase func() error) {

	log.Printf("Running test %s", name)
	err := testCase()
	if err != nil {
		log.Fatalf("Test %s failed: %s", name, err)
	}

	report.Total++
}

func (report *Report) ExpectError(name string, testCase func() error) {

	log.Printf("Running test %s", name)
	err := testCase()
	if err == nil {
		log.Fatalf("Test %s failed: expected error", name)
	}

	report.Total++
}
