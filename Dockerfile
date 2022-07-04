#
# RestGoMail Dockerfile
#

FROM golang AS restgomailbuildstage
RUN mkdir /restgomail
COPY *.go /restgomail/
ENV GO111MODULE=auto
RUN cd /restgomail \
    && go get github.com/hyper-prog/smartjson \
    && CGO_ENABLED=0 GOOS=linux go build -a -o restgomail restgomail.go smtpassembler.go

FROM alpine AS restgomail
LABEL maintainer="hyper80@gmail.com" Description="RestGoMail - HTTPS REST-Go-MAIL (SMTP) gateway"
COPY --from=restgomailbuildstage /restgomail/restgomail /usr/local/bin
RUN mkdir /restgomail
VOLUME ["/restgomail"]
WORKDIR /restgomail
CMD ["/usr/local/bin/restgomail","server.json"]
