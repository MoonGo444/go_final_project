FROM golang:1.24-alpine AS builder

WORKDIR /myapp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o goapp .

FROM ubuntu:latest

WORKDIR /myapp

COPY --from=builder /myapp/goapp /myapp/goapp
COPY --from=builder /myapp/web/ /myapp/web/

CMD ["/myapp/goapp"]