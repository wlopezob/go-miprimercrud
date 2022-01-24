FROM golang:alpine AS builder
RUN apt-get update
ENV GOO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build src/main.go
CMD [ "./main" ]

# Building image with the binary
#FROM scratch
#COPY --from=builder /go/src/app .
#EXPOSE 8080
#ENTRYPOINT ["/main"]
