package token

type Type string

const (
	EmailVerification Type = "email_verification"
	PasswordReset     Type = "password_reset"
	MagicLink         Type = "magic_link"
	EmailChange       Type = "email_change"
)
