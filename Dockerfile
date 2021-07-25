FROM golang:1.16-alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shoppinglist-backend .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder ["/build/shoppinglist-backend", "/build/.env", "/"]

ENTRYPOINT ["/shoppinglist-backend"]

