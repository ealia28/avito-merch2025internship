FROM golang:1.27 AS builder

WORKDIR ${GOPATH}/src
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /build ./cmd/shop
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build /shop
EXPOSE 8080
ENTRYPOINT ["/shop"]