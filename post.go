package zif

import "encoding/json"

const (
	TitleMax    = 144
	TagsMax     = 256
	MaxPostSize = TitleMax + TagsMax + 1024
)

type Post struct {
	Id         int
	InfoHash   string
	Title      string
	Size       int
	FileCount  int
	Seeders    int
	Leechers   int
	UploadDate int
	Tags       string
}

func (p Post) Json() ([]byte, error) {
	json, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	return json, nil
}
