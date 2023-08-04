package data

import "io"

type File struct {
	FileName  string
	Extension string
	Mime      string
	Reader    io.Reader
}
