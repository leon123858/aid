# Chat Protocol

## Infrastructure

build

```bash
docker build -t ai-front .
docker run -d -p 80:80 ai-front
```

deploy

```bash
docker buildx build --platform linux/amd64 -t ai-front .
docker tag ai-front leon1234858/ai-front
docker push leon1234858/ai-front
```