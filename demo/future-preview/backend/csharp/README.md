# Demo Backend

## How to run

deploy image

```bash
#docker buildx build -f ./Dockerfile --platform linux/amd64 .
docker build -t ai-back .
docker tag ai-back leon1234858/ai-back:latest
docker push leon1234858/ai-back:latest
```

run local

```bash
docek build -t ai-back .
docker run -e APIKey="<open-ai-APIKey>" -e Model="gpt-3.5-turbo" -d -p 3000:80 ai-back 
```