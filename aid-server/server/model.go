package server

type Space struct {
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type Request struct {
	Space
}

type Response struct {
	Result  bool   `json:"result"`
	Content string `json:"content"`
}

type LoginRequest struct {
	AID  string `json:"aid"`
	Sign string `json:"sign"`
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

type TriggerRequest struct {
	Request
}

func generateResponse(result bool, content string) Response {
	return Response{
		Result:  result,
		Content: content,
	}
}
