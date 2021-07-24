FROM golang:1.16-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shoppinglist-backend .

FROM scratch

COPY --from=builder ["/build/shoppinglist-backend", "/build/.env", "/"]

ENTRYPOINT ["/shoppinglist-backend"]

