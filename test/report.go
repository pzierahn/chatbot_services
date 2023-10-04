package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type Report struct {
	Failed []string
	Passed []string
	Total  int
}

type testing struct {
	name string
}

func (t testing) pass() (pass bool) {
	var ok bool
	_, path, line, ok := runtime.Caller(1)
	if !ok {
		path = "???"
		line = 0
	}

	filename := filepath.Base(path)

	fmt.Printf("%v %s:%d: %v %v\n", time.Now().Format(time.DateTime), filename, line, t.name, "OK")

	return true
}

func (t testing) expectError(err error) (pass bool) {
	var ok bool
	_, path, line, ok := runtime.Caller(1)
	if !ok {
		path = "???"
		line = 0
	}

	filename := filepath.Base(path)

	if err != nil {
		fmt.Printf("%v %s:%d: %v %v\n", time.Now().Format(time.DateTime), filename, line, t.name, "OK")
	} else {
		fmt.Printf("%v %s:%d: %v %v\n", time.Now().Format(time.DateTime), filename, line, t.name, "Expected error")
	}

	return err != nil
}

func (t testing) fail(err error) (pass bool) {
	var ok bool
	_, path, line, ok := runtime.Caller(1)
	if !ok {
		path = "???"
		line = 0
	}

	filename := filepath.Base(path)

	if err == nil {
		fmt.Printf("%v %s:%d: %v %v\n", time.Now().Format(time.DateTime), filename, line, t.name, "OK")
	} else {
		fmt.Printf("%v %s:%d: %v %v\n", time.Now().Format(time.DateTime), filename, line, t.name, err)
	}

	return err == nil
}

func (report *Report) Run(name string, testCase func(testing) bool) {

	passed := testCase(testing{name: name})
	if passed {
		report.Passed = append(report.Passed, name)
	} else {
		report.Failed = append(report.Failed, name)
	}

	report.Total++
}
