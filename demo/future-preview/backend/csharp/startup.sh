#!/bin/bash

cd backend/csharp

# Set environment variables
#read -p "Use Azure OpenAI? (true/false): " UseAzureOpenAI
UseAzureOpenAI=false
export UseAzureOpenAI
#if [ "$UseAzureOpenAI" = "true" ]; then
#    read -p "Enter Azure Deployment: " AzureDeployment
#    read -p "Enter Azure Endpoint: " AzureEndpoint
#    export AzureDeployment
#    export AzureEndpoint
#    echo "Please sign in to Azure:"
#    az login
#else
read -p "Enter OpenAI API Key (default is <fake key>): " APIKey
read -p "Enter OpenAI Model (default is gpt-3.5-turbo): " Model
APIKey=${APIKey:-<fake key>}
Model=${Model:-gpt-3.5-turbo}
export APIKey
export Model
#fi

# Restore dependencies
echo "Restoring dependencies..."
dotnet restore

# Run the application
echo "Running the application..."
dotnet run