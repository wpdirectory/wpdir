# Frontend Build
FROM node:latest AS node-env
ADD web/. /web/
WORKDIR /web/
RUN npm install

# Build Stage
FROM golang:1.10.3 AS go-env
ADD . /go/src/github.com/wpdirectory/wpdir
RUN cd /go/src/github.com/wpdirectory/wpdir && go get -d -v

# Embed Static Files Into Go
COPY --from=node-env /web /go/src/github.com/wpdirectory/wpdir
WORKDIR /go/src/github.com/wpdirectory/wpdir/scripts/assets/
RUN go get -d -v && go build assets_generate.go

# Compile Binary
WORKDIR /go/src/github.com/wpdirectory/wpdir
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wpdir .

# Final Stage
FROM alpine:latest
LABEL maintainer="Peter Booker <mail@peterbooker.com>"

RUN apk --no-cache add ca-certificates
COPY --from=go-env /go/src/github.com/wpdirectory/wpdir/wpdir /usr/local/bin
WORKDIR /etc/wpdir

ENTRYPOINT ["wpdir"]