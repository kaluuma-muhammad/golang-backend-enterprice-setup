package email

import (
	"bytes"
	"embed"
	"fmt"
	htmlTemplate "html/template"
)

//go:embed templates/*.html
var htmlFS embed.FS

type Renderer struct {
	templates *htmlTemplate.Template
}

func NewRenderer() (*Renderer, error) {
	t, err := htmlTemplate.ParseFS(htmlFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &Renderer{
		templates: t,
	}, nil
}

func (r *Renderer) RenderHTML(name string, data any) (string, error) {
	t, err := htmlTemplate.ParseFS(htmlFS, "templates/"+name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
