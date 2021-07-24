FROM golang:1.16-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o shoppinglist-backend .

FROM scratch

COPY --from=builder ["/build/shoppinglist-backend", "/build/.env", "/"]

ENTRYPOINT ["/shoppinglist-backend"]

