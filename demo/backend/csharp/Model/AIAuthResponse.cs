namespace Backend.Model;

public class AIAuthResponse
{
    public string? Token { get; private set; }
    public string? Uuid { get; private set; }
    public string? Message { get; private set; }
    public bool? Result { get; private set; }

    private AIAuthResponse() { }

    public class Builder
    {
        private string? _token;
        private string? _uuid;
        private string? _message;
        private bool? _result;

        public Builder WithToken(string? token)
        {
            _token = token;
            return this;
        }

        public Builder WithUuid(string? uuid)
        {
            _uuid = uuid;
            return this;
        }

        public Builder WithMessage(string? message)
        {
            _message = message;
            return this;
        }

        public Builder WithResult(bool? result)
        {
            _result = result;
            return this;
        }

        public AIAuthResponse Build()
        {
            return new AIAuthResponse
            {
                Token = _token,
                Uuid = _uuid,
                Message = _message,
                Result = _result
            };
        }
    }
}