package data

import (
	"bytes"
	"io"
	"strconv"
)

// Convert a post to a string, with an optional separator between fields, and
// an optional terminating value (appended to the end of the post string).
// This is *actually* signed, to allow for different encoders encoding in a
// different order. (relying on a json encoding to always encode the same way
// is not the best idea)
func PostToString(p *Post, sep, term string) string {
	buf := bytes.Buffer{}

	WritePost(p, sep, term, &buf)

	return buf.String()
}

func WritePost(p *Post, sep, term string, w io.Writer) {
	w.Write([]byte(strconv.Itoa(p.Id)))
	w.Write([]byte(sep))
	w.Write([]byte(p.InfoHash))
	w.Write([]byte(sep))
	w.Write([]byte(p.Title))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.Size)))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.FileCount)))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.Seeders)))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.Leechers)))
	w.Write([]byte(sep))
	w.Write([]byte(strconv.Itoa(p.UploadDate)))
	w.Write([]byte(sep))
	w.Write([]byte(p.Tags))
	w.Write([]byte(sep))
	w.Write([]byte(p.Meta))
	w.Write([]byte(sep))
	w.Write([]byte(term))
}
