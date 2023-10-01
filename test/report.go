package test

import "log"

type Report struct {
	Total int
}

func (report *Report) Run(name string, testCase func() error) {

	err := testCase()
	if err != nil {
		log.Printf("Test %s failed: %s", name, err)
	} else {
		log.Printf("Passsing %s test", name)
	}

	report.Total++
}

func (report *Report) ExpectError(name string, testCase func() error) {

	log.Printf("Running test %s", name)
	err := testCase()
	if err == nil {
		log.Printf("Test %s failed: expected error", name)
	} else {
		log.Printf("Passsing %s test", name)
	}

	report.Total++
}
