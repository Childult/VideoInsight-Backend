package job

import (
	"fmt"
	"swc/mongodb"
)

// Job include media and audio
type Job struct {
	JobID    string   `bson:"job_id"              json:"job_id"`
	DeviceID string   `bson:"device_id"           json:"device_id"`
	URL      string   `bson:"url"                 json:"url"`
	KeyWords []string `bson:"key_words"           json:"key_words"`
	Status   int32    `bson:"status"              json:"status"`
	AbsText  string   `bson:"abstract_text"       json:"abstract_text"`
	AbsVideo string   `bson:"abstract_video"      json:"abstract_video"`
}

// GetKeyTag implement the interface Key
func (j Job) GetKeyTag() string {
	return "job_id"
}

// GetKeyValue implement the interface Key
func (j Job) GetKeyValue() string {
	return j.JobID
}

// GetCollName implement the interface Key
func (j Job) GetCollName() string {
	return "job"
}

// SetStatus implement the interface Key
func (j *Job) SetStatus(i int32) {
	j.Status = i
	err := mongodb.Update(j)
	fmt.Println(err)
}

// SetAbsText implement the interface Key
func (j *Job) SetAbsText(key string) {
	j.AbsText = key
	mongodb.Update(j)
}

// SetAbsVideo implement the interface Key
func (j *Job) SetAbsVideo(key string) {
	j.AbsVideo = key
	mongodb.Update(j)
}
