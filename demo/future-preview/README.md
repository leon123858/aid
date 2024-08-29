# Microsoft AI Chat Protocol Samples

This is copy from `https://github.com/microsoft/ai-chat-protocol`

The goal of this project is want to preview how will aid be used in the future.

This domo do not include real aid system module, if you want to find real aid system, please visit other demo.

## Environment

1. Node.js
2. .NET 7.0

## Frontend

1. Clone the repository to your machine.
2. In one terminal, navigate to the `frontend/js/react` directory.
3. In the `frontend/js/react` directory, run `npm install` to install your dependencies, including [`@microsoft/ai-chat-protocol`](https://www.npmjs.com/package/@microsoft/ai-chat-protocol).
4. Next, run `npm run dev` to start your web application.

## Backend

1. create a OpenAI API Key in `https://platform.openai.com/api-keys`
2. Clone the repository to your machine.
3. In one terminal, navigate to the `backend/charp` directory.
4. run `bash ./startup.sh` to start the backend server.
5. you should paste your OpenAI API Key in the terminal.

## Demo Script

this step depend on local ip setting
- wallet/lib/constants/config.dart
- demo/frontend/js/react/vite.config.ts
- demo/backend/csharp/Services/AID.cs

1. Open App in MacOS
    ```bash
    cd ./wallet/build/macos/Build/Products/Release
    # use `open` to open the app
    open wallet.app
    ```
2. Open App in Browser
    ```bash
    cd ./demo/frontend/js/react
    npm run dev
    ```
3. Open aid server
    ```bash
    cd ./aid-server
    make
    ```
4. Open App backend server
    ```bash
    cd ./demo/backend/csharp
    bash ./startup.sh
    ```
