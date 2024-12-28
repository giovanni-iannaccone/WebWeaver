FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o webweaver ./cmd/main.go

CMD ["./webweaver", "-c", "configs/configs.json"]
