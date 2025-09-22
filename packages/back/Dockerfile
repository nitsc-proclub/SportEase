FROM golang:1.24-alpine AS builder

# go install でインストールしたコマンドを使えるようにPATHを通す
ENV PATH="/go/bin:${PATH}"

WORKDIR /app

# swag CLIをインストール
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN swag init -g ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /root/

COPY --from=builder /go/bin/app .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./app"]
