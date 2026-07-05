package token

type Type string

const (
	AccountVerification Type = "account_verification"
	EmailVerification   Type = "email_verification"
	PasswordReset       Type = "password_reset"
	MagicLink           Type = "magic_link"
	EmailChange         Type = "email_change"
)
