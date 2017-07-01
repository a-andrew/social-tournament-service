FROM golang:1.7.1-alpine

RUN apk add git --no-cache && \
    go get -u github.com/kardianos/govendor
      
ADD . /go/src/github.com/social-tournament-service/

RUN cd /go/src/github.com/social-tournament-service/app && \
    govendor sync && \
    go build -o /go/bin/social-tournament-service && \
    apk del git


WORKDIR /go/bin/

CMD /go/bin/social-tournament-service

EXPOSE 8080