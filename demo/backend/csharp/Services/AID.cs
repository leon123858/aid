using System.Text;
using System.Text.Json;
using System.Net.Http.Headers;

namespace Backend.Services;

public class AIDService
{
    private readonly HttpClient client = new HttpClient();
    private readonly string _aidUrl = "";
    
    public AIDService()
    {
        client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
        _aidUrl = Environment.GetEnvironmentVariable("AID_URL") ?? "http://localhost:8080/";
    }
    
    private async Task<AIDResponse> PostToGoService(string endpoint, object request)
    {
        var json = JsonSerializer.Serialize(request);
        var content = new StringContent(json, Encoding.UTF8, "application/json");
        var response = await client.PostAsync($"{_aidUrl}api/{endpoint}", content);
        var res = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<AIDResponse>(res);
    }
    
    public async Task<AIDResponse> ask(string ip, string browser)
    {
        var request = new AskRequest(ip, browser);
        var response = await PostToGoService("ask", request);
        if (response.result)
        {
            return response;
        }
        throw new Exception("error: " + response.content);
    }
    
    public async Task<AIDResponse> check(string uid, string ip, string browser)
    {
        var request = new CheckRequest(uid, ip, browser);
        var response = await PostToGoService("check", request);
        if (response.result)
        {
            return response;
        }
        throw new Exception("error: " + response.content);
    }
    
    public async Task<AIDResponse> verify(string token, string uid)
    {
        var request = new VerifyRequest(uid);
        client.DefaultRequestHeaders.Remove("Authorization");
        client.DefaultRequestHeaders.Add("Authorization", token); 
        var response = await PostToGoService("verify", request);
        if (response.result)
        {
            return response;
        }
        throw new Exception("error: " + response.content);
    }
    
    private class AskRequest
    {
        public string IP { get; set; }
        public string Browser { get; set; }
        
        public AskRequest(string ip, string browser)
        {
            IP = ip;
            Browser = browser;
        }
    }

    private class CheckRequest
    {
        public string UID { get; set; }
        public string IP { get; set; }
        public string Browser { get; set; }
        
        public CheckRequest(string uid, string ip, string browser)
        {
            UID = uid;
            IP = ip;
            Browser = browser;
        }
    }

    private class VerifyRequest
    {
        public string UID { get; set; }
        
        public VerifyRequest(string uid)
        {
            UID = uid;
        }
    }

    public class AIDResponse
    {
        public  bool result { get; set; }
        public  string content { get; set; }
    }
}

