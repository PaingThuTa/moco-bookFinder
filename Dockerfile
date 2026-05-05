FROM alpine:latest AS builder
RUN apk add --no-cache git curl gcc musl-dev
ENV GOVERSION=1.25.7
RUN curl -fsSL "https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:$PATH"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bot .
EXPOSE 10000
CMD ["./bot"]
