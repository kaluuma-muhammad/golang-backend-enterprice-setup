package email

import "io"

type Message struct {
	To          string
	Subject     string
	HTML        string
	CC          []string
	BCC         []string
	ReplyTo     string
	Attachments []Attachment
}

type Attachment struct {
	Filename string
	Reader   io.Reader
}
