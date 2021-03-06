#
# RestGoMail Dockerfile
#

FROM golang AS restgomailbuildstage
RUN mkdir /restgomail
COPY restgomail.go smartjson.go /restgomail/
RUN cd /restgomail && CGO_ENABLED=0 GOOS=linux go build -a -o restgomail restgomail.go smartjson.go

FROM alpine AS restgomail
LABEL maintainer="hyper80@gmail.com" Description="RestGoMail - HTTPS REST-Go-MAIL (SMTP) gateway"
COPY --from=restgomailbuildstage /restgomail/restgomail /usr/local/bin
RUN mkdir /restgomail
VOLUME ["/restgomail"]
WORKDIR /restgomail
CMD ["/usr/local/bin/restgomail","server.json"]
