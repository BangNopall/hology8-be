FROM golang:1.25.8-alpine AS builder 

RUN apk update && apk add --no-cache git 

WORKDIR /app    

COPY go.mod go.sum ./ 

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app 

FROM alpine:3.20

WORKDIR /app 

COPY --from=builder /app/main . 
COPY --from=builder /app/data ./data 
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env .

EXPOSE 8080

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT [ "/app/main" ]
