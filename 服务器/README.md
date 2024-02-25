## 服务器

```bash
docker run -d --restart always --name 2024challenge -v "$(readlink -f ~/go)":/go -v "$(pwd)":/app -w /app -p 3000:3000 golang:1.22 go run .
```
