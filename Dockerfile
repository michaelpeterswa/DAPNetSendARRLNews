FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go get -d -v .
RUN go build .

CMD ./DAPNetSendARRLNews