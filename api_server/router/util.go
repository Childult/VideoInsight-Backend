package router

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// MessageJSON contains unique identity and user data
type MessageJSON struct {
	DeviceID string   `json:"deviceid"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
}

func (json MessageJSON) String() string {
	return json.DeviceID + json.URL + strings.Join(json.KeyWords, "")
}

// GetHash will return a hash
func (json MessageJSON) GetHash() (result [12]byte) {
	hash := sha1.New()
	hash.Write([]byte(json.String()))
	copy(result[:], hash.Sum([]byte(""))[0:12])
	return
}

// GetID will return a id
func (json MessageJSON) GetID() string {
	return fmt.Sprintf("%v", json.GetHash())
}

func getJSON(c *gin.Context) (json MessageJSON, err error) {
	// 获取数据
	err = c.ShouldBindJSON(&json)
	if err != nil {
		err = fmt.Errorf("%s", gin.H{"error": "Wrong Format"})
		return
	}
	json.KeyWords = removeEmptyString(json.KeyWords)
	return
}

func removeEmptyString(a []string) []string {
	return deleteKeywords(a, "")
}

func deleteKeywords(a []string, s string) []string {
	j := 0
	for _, val := range a {
		if val == s {
			a[j] = val
			j++
		}
	}
	return a[j:]
}
