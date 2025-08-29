package test

import "flag"

func IsTestRun() bool {
	return flag.Lookup("test.v") != nil
}
