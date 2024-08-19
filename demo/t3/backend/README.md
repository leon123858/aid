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
        json result = {"hash", *it};
        api->writeContractState(&result.dump());
      }
      // empty
      json emptyResult = {"hash", ""};
      api->writeContractState(&emptyResult.dump());
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