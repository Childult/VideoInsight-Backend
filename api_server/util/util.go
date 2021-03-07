package util

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// WorkSpace 动态设置工作路径
	WorkSpace = ""
	// SavePath 文件保存位置
	SavePath = "/home/donwload"
)

// SetWorkSpace 获取当前路径
func SetWorkSpace() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	WorkSpace, _ = filepath.Split(dir)
}

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

// GetJSON return a json
func GetJSON(c *gin.Context) (json MessageJSON, err error) {
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
