FROM golang:1.23 AS builder

RUN apt-get update && apt-get install -y tzdata bash
WORKDIR /web
COPY . .
#ENV CSG_ENABLED=0
#ENV GOOS=linux
#ENV GOARCH=amd64
RUN go mod download
ENV TZ=America/Toronto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:3.20
RUN apk update && apk add --no-cache tzdata bash
WORKDIR /srv
ENV TZ=America/Toronto
COPY --from=builder  /web/app .
#COPY --from=builder  /web/templates ./templates

EXPOSE 80
CMD [ "./app" ]