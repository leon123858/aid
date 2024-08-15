using Backend.Model;
using Backend.Services;
using Microsoft.AspNetCore.Mvc;

namespace Backend.Controllers;

[ApiController]
[Route("api/[controller]")]
public class AuthController : ControllerBase
{
    private readonly UserDatabase _userDatabase;
    
    public AuthController(UserDatabase userDatabase)
    {
        _userDatabase = userDatabase;
    }
    
    [HttpPost("login")]
    [Consumes("application/json")]
    public async Task<IActionResult> Login(AIAuthRequest request)
    {
        // try to find user in database
        try
        {
            var uidList = _userDatabase.ValidateUser(request.Username, request.Password);
            // process user list
            Console.WriteLine("db:" + uidList.Count);
            foreach (var uid in uidList)
            { 
                Console.WriteLine(uid);   
            }
            switch (uidList.Count)
            {
                case > 1:
                    // for all users call AID check method
                    var service = new AIDService();
                    var onlineList = new List<string>();
                    if (request.Token != "")
                    {
                        foreach (var uid in uidList)
                        {
                            try
                            {
                                var res = await service.verify(request.Token, uid);
                                if (res.result)
                                {
                                    onlineList.Add(uid);
                                }
                            }
                            catch (Exception e)
                            {
                                Console.WriteLine(e.Message);
                            }

                        }
                    }
                    else
                    {
                        foreach (var uid in uidList)
                        {
                            try
                            {
                                var res = await service.check(uid, request.IP, request.fingerprint);
                                Console.WriteLine(res.result + "/" + res.content);
                                if (res.result & res.content == "online")
                                {
                                    onlineList.Add(uid);
                                }
                            } catch (Exception e)
                            {
                                Console.WriteLine(e.Message);
                            }
                        }
                    }
                    Console.WriteLine("online:" + onlineList.Count);
                    foreach (var uid in onlineList)
                    { 
                        Console.WriteLine(uid);   
                    }
                    return onlineList.Count switch
                    {
                        1 => Ok(new AIAuthResponse.Builder().WithUuid(onlineList[0])
                            .WithMessage("Successfully logged in")
                            .WithResult(true)
                            .Build()),
                        0 => BadRequest(new AIAuthResponse.Builder().WithMessage("No user found")
                            .WithResult(false)
                            .Build()),
                        _ => BadRequest(new AIAuthResponse.Builder().WithMessage("Multiple users found")
                            .WithResult(false)
                            .Build())
                    };
                case 1:
                    return Ok(new AIAuthResponse.Builder()
                        .WithUuid(uidList[0])
                        .WithMessage("Successfully logged in")
                        .WithResult(true)
                        .Build());
                case 0:
                    return BadRequest(new AIAuthResponse.Builder()
                        .WithMessage("no user recognized")
                        .WithResult(false)
                        .Build());
                default:
                    return BadRequest(new AIAuthResponse.Builder()
                        .WithMessage("internal error")
                        .WithResult(false)
                        .Build());
            }
            
        } catch (Exception e)
        {
            return BadRequest(new AIAuthResponse.Builder()
                .WithMessage(e.Message)
                .WithResult(false)
                .Build());
        }
    }

    [HttpPost("register")]
    [Consumes("application/json")]
    public async Task<IActionResult> Register(AIAuthRequest request)
    {
        var service = new AIDService();
        var resBuilder = new AIAuthResponse.Builder();
        // use AID ask method to get new uid
        try
        {
            var res = await service.ask(request.IP, request.fingerprint);
            var uid = res.content;
            // save new user to database
            _userDatabase.AddUser(uid, request.Username, request.Password);
            return Ok(resBuilder
                .WithUuid(uid)
                .WithMessage("Successfully registered")
                .WithResult(true)
                .Build());
        } catch (Exception e)
        {
            return BadRequest(resBuilder
                .WithMessage(e.Message)
                .WithResult(false)
                .Build());
        }
    }
}