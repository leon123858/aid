FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 80

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["ChatProtocolBackend.csproj", "./"]
# disable signature verification to speed up nuget restore
ENV DOTNET_NUGET_SIGNATURE_VERIFICATION=false
RUN dotnet restore "ChatProtocolBackend.csproj"
COPY . .
WORKDIR "/src/"
RUN dotnet build "ChatProtocolBackend.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "ChatProtocolBackend.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .

# should use -e APIKey="$APIKey" -e Model="$Model" set environment variables outside of docker
ENTRYPOINT ["dotnet", "ChatProtocolBackend.dll"]
