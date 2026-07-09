package loginhistory

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusLogout  Status = "logout"
	StatusExpired Status = "expired"
	StatusRevoked Status = "revoked"
)
