FROM golang:1.19.3-alpine3.16

COPY . /go/src/myapp
WORKDIR /go/src/myapp/cmd/simple-S3-cache
RUN go mod download
RUN go build -o /simple-S3-cache
EXPOSE 8084
CMD [ "/simple-S3-cache" ]