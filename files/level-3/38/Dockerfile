FROM alpine:3

RUN mkdir /c2
COPY ./server.go /c2
RUN apk add go curl vim

RUN go run /c2/server.go -p 80
