package storage

import "io"

type UploadInput struct {
	File        io.Reader
	Filename    string
	Folder      string
	ContentType string
}

type File struct {
	Path        string
	URL         string
	Filename    string
	Size        int64
	ContentType string
}
