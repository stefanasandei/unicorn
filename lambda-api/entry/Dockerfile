FROM golang:1.20

COPY . $GOPATH/app

WORKDIR $GOPATH/app/entry

RUN go get .

RUN go build .

EXPOSE 3000

CMD ./entry
