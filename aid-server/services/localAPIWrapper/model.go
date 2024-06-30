package localAPIWrapper

type UsageRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Token       string `json:"token,omitempty"`
	IP          string `json:"ip,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type UsageResponse struct {
	Token   string `json:"token,omitempty"`
	UUID    string `json:"uuid,omitempty"`
	Message string `json:"message,omitempty"`
	Result  bool   `json:"result"`
}
