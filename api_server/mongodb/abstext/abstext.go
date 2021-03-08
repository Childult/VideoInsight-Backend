package abstext

// AbsText 文本摘要
type AbsText struct {
	Hash     string   `bson:"hash"                json:"hash"`
	URL      string   `bson:"url"                 json:"url"`
	KeyWords []string `bson:"key_words,omitempty" json:"key_words,omitempty"`
	Text     string   `bson:"text"                json:"text"`
	Abstract string   `bson:"abstract"            json:"abstract"`
}

// GetKeyTag implement the interface Key
func (media AbsText) GetKeyTag() string {
	return "hash"
}

// GetKeyValue implement the interface Key
func (media AbsText) GetKeyValue() string {
	return media.Hash
}

// GetCollName implement the interface Key
func (media AbsText) GetCollName() string {
	return "abstract_text"
}
