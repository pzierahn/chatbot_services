package test

import "log"

type Report struct {
	Total int
}

func (report *Report) Run(name string, testCase func() error) {

	err := testCase()
	if err != nil {
		log.Printf("Failed %s failed: %s", name, err)
	} else {
		log.Printf("Passsing %s test", name)
	}

	report.Total++
}

func (report *Report) ExpectError(name string, testCase func() error) {

	err := testCase()
	if err == nil {
		log.Printf("Failed %s failed: expected error", name)
	} else {
		log.Printf("Passsing %s test", name)
	}

	report.Total++
}
