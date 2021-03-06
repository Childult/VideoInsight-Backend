package server

import (
	"swc/mongodb"
)

// StartTask .
func StartTask(job mongodb.Job) {
	err := downloadMedia(&job)
	if err == completed {

	} else if err == exist {

	}
}
