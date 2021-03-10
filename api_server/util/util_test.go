package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteKeywords(t *testing.T) {
	// 测试函数 deleteKeywords
	tests := []struct {
		rawSlice     []string
		targetStr    string
		expectResult []string
	}{
		{nil, "", []string{}},
		{[]string{}, "", []string{}},
		{[]string{""}, "", []string{}},
		{[]string{"hello", "", "world"}, "", []string{"hello", "world"}},
	}

	for _, test := range tests {
		result := deleteKeywords(test.rawSlice, test.targetStr)
		assert.Equal(t, test.expectResult, result)
	}
}
