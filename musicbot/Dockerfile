FROM golang:alpine3.15 AS BUILD
WORKDIR /api
COPY go.sum .
COPY go.mod .
COPY main.go .
RUN go install
RUN go build -o api main.go

FROM alpine:3.15
ENV JAR_FILE=JMusicBot-0.3.8.jar

RUN  apk update \
  && apk upgrade \
  && apk add ca-certificates \
  && update-ca-certificates \
  && apk add --update coreutils && rm -rf /var/cache/apk/*   \ 
  && apk add --update openjdk11 tzdata curl unzip bash \
  && apk add --no-cache nss \
  && rm -rf /var/cache/apk/*

WORKDIR /jmusicbot
COPY ./${JAR_FILE} .
COPY --from=BUILD /api/api .
RUN chmod +x api

CMD [ "/jmusicbot/api"]

EXPOSE 8080