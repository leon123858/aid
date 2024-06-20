package res

type Response struct {
	Result  bool   `json:"result"`
	Content string `json:"content"`
}

func GenerateResponse(result bool, content string) Response {
	return Response{
		Result:  result,
		Content: content,
	}
}
