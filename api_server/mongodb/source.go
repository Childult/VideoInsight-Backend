package mongodb

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
