package util

import (
	"fmt"
	"testing"
)

func TestGetJobStatus(t *testing.T) {
	var x int32 = 1
	for x <= JobErrVideoAnalysisGRPCallJobIDNotMatch {
		fmt.Printf("%d:%s\n", x, GetJobStatus(x))
		x <<= 1
	}
}
