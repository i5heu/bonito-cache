FROM golang:1.19.3-alpine3.16

COPY . /go/src/myapp
WORKDIR /go/src/myapp
RUN go mod download
RUN go build -o /godocker
EXPOSE 8084
CMD [ "/godocker" ]