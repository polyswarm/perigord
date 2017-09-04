package main

import (
	"testing"

	_ "{{.project}}/tests"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner
func Test(t *testing.T) { TestingT(t) }
