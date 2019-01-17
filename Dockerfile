FROM golang:1.11.3

MAINTAINER Acke vimer757216574@gmail.com

WORKDIR /go/src/awesomeProject/src
COPY . /go/src/awesomeProject/

RUN go get github.com/facebookgo/inject
RUN go get github.com/jinzhu/gorm
#RUN go get gopkg.in/mgo.v2/bson

RUN go build .

EXPOSE 12019
ENTRYPOINT ["./src"]
