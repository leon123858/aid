# aid
this project include aid server implementation and full stack demo.

## Concept

provide rest api for service to check user identity message.

## paper

- [Design and Implementation of Autonomous Identity System Based on OurChain](https://github.com/leon123858/aid-paper)

## aid system

![img.png](doc/overview-dark.png)

- OurChain is a blockchain system that can sync data between different aid-server.
- Wallet can ask OurChain to do some operation.
- Wallet is embedded with frontend application, which can communicate with service.
- Service is a backend application that can communicate with aid-server.

related project:
- [aid-server](https://github.com/leon123858/aid)
- [aid-wallet](https://github.com/leon123858/aidjs)
- [aid-service](https://github.com/leon123858/aidgo)
- [OurChain](https://github.com/leon123858/OurChain)
- [OurChain-Agent](https://github.com/leon123858/ourchain-agent)
