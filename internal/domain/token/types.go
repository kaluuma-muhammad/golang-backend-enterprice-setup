package token

type Type string

const (
	AccountActivation  Type = "account_activation"
	PasswordReset      Type = "password_reset"
	PasswordResetGrant Type = "password_reset_grant"
	MagicLink          Type = "magic_link"
	EmailChange        Type = "email_change"
)
