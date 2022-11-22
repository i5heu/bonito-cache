FROM golang:1.19.3-alpine3.16

COPY . /go/src/myapp
WORKDIR /go/src/myapp/cmd/bonito-cache
RUN go mod download
RUN go build -o /bonito-cache
EXPOSE 8084
CMD [ "/bonito-cache" ]