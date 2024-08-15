package server

type Space struct {
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type Request struct {
	Space
}

type LoginRequest struct {
	AID       string `json:"aid"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Request
}

type RegisterRequest struct {
	AID       string `json:"aid"`
	PublicKey string `json:"publicKey"`
	Request
}

type LogoutRequest struct {
	AID string `json:"aid"`
	Request
}

type AskRequest struct {
	Request
}

type CheckRequest struct {
	UID string `json:"uid"`
	Request
}

type VerifyRequest struct {
	UID string `json:"uid"`
	Request
}
