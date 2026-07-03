package email

import "fmt"

func VerificationTemplate(name string, code string) (string, string, string) {
	subject := "Verify your email"

	html := fmt.Sprintf(`
		<h2>Hello %s</h2>

		<p>Your verification code is</p>
		<h1>%s</h1>
		<p>This code expires in 15 minutes.</p>
	`, name, code)

	text := fmt.Sprintf(
		"Hello %s\n\nYour verification code is %s",
		name,
		code,
	)

	return subject, html, text
}
