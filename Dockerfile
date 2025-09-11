FROM golang:1.21

WORKDIR /app
COPY backend/ .

# Run module customization script
ENV MODULE_NAME=github.com/alimosavifard/zyros-backend
RUN chmod +x customize-module.sh && ./customize-module.sh

RUN go mod tidy
COPY backend/.env .

EXPOSE 8080
CMD ["go", "run", "main.go"]