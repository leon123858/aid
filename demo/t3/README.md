# DEMO of Type3

this is a demo of type3 AID

## Concept

a todo list can use aid to log in a remote todo list application, which can share the todo list with others.

![image](../../doc/overview-t3.png)

## Smart Contract

basic smart contract code(only for demo):

```cpp
#include <iostream>
#include <json.hpp>
#include "ourcontract.h"

using json = nlohmann::json;
using namespace std;

// contract main function
extern "C" int contract_main(void *arg)
{
  // cast argument
  ContractArguments *contractArg = (ContractArguments *)arg;
  ContractAPI *api = &contractArg->api;
  // pure call contract
  if (contractArg->isPureCall)
  {
    string command = contractArg->parameters[0];
    api->contractLog("command: " + command);
    if (command == "get")
    {
      api->generalContractInterfaceOutput("aid", "0.1.0");
      return 0;
    }
    else if (command == "verify")
    {
      // pure operation
      string state = api->readContractState();
      json j = json::parse(state);
      auto it = j.find(contractArg->parameters[1]);
      if (it != j.end())
      {
        json result = json::object();
        result["hash"] = it.value();
        string str = result.dump();
        api->writeContractState(&str);
        return 0;
      }
      // empty
      json result = json::object();
      result["hash"] = "";
      string str = result.dump();
      api->writeContractState(&str);
      return 0;
    }
  }
  // non-pure call contract
  string state = api->readContractState();
  // deploy contract init call
  if (state == "null")
  {
    json j = json::object();
    // write contract state
    state = j.dump();
    api->writeContractState(&state);
    return 0;
  }
  if (contractArg->parameters[0] == "register")
  {
    json j = json::parse(state);
    j[contractArg->parameters[1]] = contractArg->parameters[2];
    // write contract state
    state = j.dump();
    api->writeContractState(&state);
    return 0;
  }
  return 0;
}
```

## How to execute

1. start the backend (use `make`)
2. start the frontend (use `yarn start`)
3. go to project `ourchain-agent` and run `docker-compose up`
4. use `web-cli` in docker to deploy the contract above
5. use `web-cli` in docker to get the contract information (wallet address,private key,contract address)
6. set above information in `frontend` project

## Note

below contract have comment function, which is more complex than above contract

if you want to enumerate the aid system more close to paper, you can use below contract

```cpp
#include <iostream>
#include <json.hpp>
#include "ourcontract.h"

using json = nlohmann::json;
using namespace std;

// contract main function
extern "C" int contract_main(void *arg)
{
  // cast argument
  ContractArguments *contractArg = (ContractArguments *)arg;
  ContractAPI *api = &contractArg->api;
  // pure call contract
  if (contractArg->isPureCall)
  {
    string command = contractArg->parameters[0];
    api->contractLog("command: " + command);
    if (command == "get")
    {
      api->generalContractInterfaceOutput("aid", "0.1.0");
      return 0;
    }
    else if (command == "verify")
    {
      // pure operation
      string state = api->readContractState();
      json j = json::parse(state);
      auto it = j.find(contractArg->parameters[1]);
      if (it != j.end())
      {
        json result = json::object();
        result["hash"] = it.value()["hash"];
        string str = result.dump();
        api->writeContractState(&str);
        return 0;
      }
      // empty
      json result = json::object();
      result["hash"] = "";
      string str = result.dump();
      api->writeContractState(&str);
      return 0;
    }
    else if (command == "lsc")
    {
      // pure operation
      string state = api->readContractState();
      json j = json::parse(state);
      auto it = j.find(contractArg->parameters[1]);
      if (it != j.end())
      {
        json result = json::object();
        result["comment"] = it.value()["comment"];
        string str = result.dump();
        api->writeContractState(&str);
        return 0;
      }
      // empty
      json result = json::object();
      result["comment"] = json::array();
      string str = result.dump();
      api->writeContractState(&str);
      return 0;
    }
    return 0;
  }
  // non-pure call contract
  string state = api->readContractState();
  // deploy contract init call
  if (state == "null")
  {
    json j = json::object();
    // write contract state
    state = j.dump();
    api->writeContractState(&state);
    return 0;
  }
  if (contractArg->parameters[0] == "register")
  {
    json j = json::parse(state);
    json element = json::object();
    element["hash"] = contractArg->parameters[2];
    element["comment"] = json::array();
    j[contractArg->parameters[1]] = element;
    // write contract state
    state = j.dump();
    api->writeContractState(&state);
    return 0;
  }
  if (contractArg->parameters[0] == "comment")
  {
    json j = json::parse(state);
    auto it = j.find(contractArg->parameters[1]);
    if (it != j.end())
    {
      it.value()["comment"].push_back(contractArg->parameters[2]);
      // write contract state
      state = j.dump();
      api->writeContractState(&state);
    }
    return 0;
  }
  return 0;
}
```