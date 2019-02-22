FROM golang:1.11.5
WORKDIR /go/src/github.com/johscheuer/event-trigger-github
RUN go get -d -v golang.org/x/net/html  
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o event-trigger-github .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/johscheuer/event-trigger-github .
CMD ["./event-trigger-github"] 
