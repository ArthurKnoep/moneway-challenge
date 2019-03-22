# build
FROM golang:latest AS builder
RUN mkdir -p /go/src/github.com/ArthurKnoep/moneway-challenge
WORKDIR /go/src/github.com/ArthurKnoep/moneway-challenge
RUN apt update && apt install --yes protobuf-compiler && rm -rf /var/lib/apt/lists/*
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/golang/dep/...
ADD . .

RUN dep ensure -v

RUN ./script/gen-protoc.sh
RUN go build -o /go/bin/account cmd/account/main.go

# release
FROM alpine
WORKDIR /app
RUN apk add libc6-compat
COPY --from=builder /go/bin/account /usr/bin
CMD ["/usr/bin/account"]
