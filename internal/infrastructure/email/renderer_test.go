package email

import (
	"testing"
)

func TestNewRenderer(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	if renderer == nil {
		t.Fatal("expected non-nil renderer")
	}

	// Test rendering HTML
	htmlData := struct {
		Title   string
		AppName string
		Name    string
		Code    string
		Expiry  int
	}{
		Title:   "Verification",
		AppName: "TestApp",
		Name:    "John",
		Code:    "123456",
		Expiry:  15,
	}

	htmlContent, err := renderer.RenderHTML("verification.html", htmlData)
	if err != nil {
		t.Errorf("failed to render HTML: %v", err)
	}
	if htmlContent == "" {
		t.Error("expected non-empty HTML content")
	}
}
