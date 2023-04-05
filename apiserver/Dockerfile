FROM golang:1.18 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

########################################################
FROM alpine:3.16

RUN apk --no-cache add ca-certificates libc6-compat 

WORKDIR /root/

COPY --from=builder /usr/local/bin/app ./

COPY ./keys/key.pem ./keys/key.pem

CMD ["./app"]
