package util

import (
	"fmt"
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

func TestHash(t *testing.T) {
	test := MessageJSON{
		URL:      "baidu.com",
		DeviceID: "1",
	}

	fmt.Println(test.GetID() == "e362abcbabcc76f7fee0c4c8")
}
