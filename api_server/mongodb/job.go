package mongodb

// Job include media and audio
type Job struct {
	DeviceID string   `json:"deviceid"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
	JobID    string   `json:"jobid"`
	Status   string   `json:"status"`
}

// 状态值
const (
	Downloading    = "Downloading"
	Processing     = "Processing"
	Completed      = "Completed"
	ErrorHappended = "ErrorHappended"
)

// GetKeyTag implement the interface Key
func (media Job) GetKeyTag() string {
	return "jobid"
}

// GetKeyValue implement the interface Key
func (media Job) GetKeyValue() string {
	return media.JobID
}

// GetCollName implement the interface Key
func (media Job) GetCollName() string {
	return "job"
}
