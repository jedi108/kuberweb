#FROM golang
#
## Fetch dependencies
#RUN go get github.com/tools/godep
#
## Add project directory to Docker image.
#ADD . /go/src/git.betfavorit.cf/vadim.tsurkov/kuberweb
#
#ENV USER jedi
#ENV HTTP_ADDR :8888
#ENV HTTP_DRAIN_INTERVAL 1s
#ENV COOKIE_SECRET qaBzlTixkx2c9S6i
#
## Replace this with actual PostgreSQL DSN.
#ENV DSN postgres://jedi@localhost:5432/kuberweb?sslmode=disable
#
#WORKDIR /go/src/git.betfavorit.cf/vadim.tsurkov/kuberweb
#
#RUN godep go build
#
#EXPOSE 8888
#CMD ./kuberweb

FROM repo.betfavorit.cf/golang:1.10-alpine

ENV GOPATH=/go/

WORKDIR /go/src/git.betfavorit.cf/backend/gateway
ADD . .

RUN go build -i .

CMD ["./gateway"]
