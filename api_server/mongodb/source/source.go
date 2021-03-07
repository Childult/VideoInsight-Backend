package source

import (
	"swc/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

// Source include media and audio
type Source struct {
	URL       string `json:"url"`
	Status    string `json:"status"`
	Location  string `json:"location"`
	VideoPath string `json:"videopath"`
	AudioPath string `json:"audiopath"`
}

// GetKeyTag implement the interface Key
func (media Source) GetKeyTag() string {
	return "url"
}

// GetKeyValue implement the interface Key
func (media Source) GetKeyValue() string {
	return media.URL
}

// GetCollName implement the interface Key
func (media Source) GetCollName() string {
	return "source"
}

// GetByKey return one source
func GetByKey(url string) (s Source, err error) {
	collName := s.GetCollName()
	key := s.GetKeyTag()
	value := url
	filter := bson.M{key: value}
	data, err := mongodb.FindOneByfilter(collName, filter)
	if err != nil {
		return
	}

	s.URL = url
	s.Status = data["status"].(string)
	s.Location = data["location"].(string)
	s.VideoPath = data["videopath"].(string)
	s.AudioPath = data["audiopath"].(string)

	return
}

// Refresh read the database, refresh variables in place
func (media *Source) Refresh() (err error) {
	*media, err = GetByKey(media.URL)
	return
}
