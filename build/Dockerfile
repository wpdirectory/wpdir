# Build Stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -i -o wpdir

# Final Stage
FROM alpine:latest
LABEL maintainer="Peter Booker <mail@peterbooker.com>"

COPY --from=build-env /src/wpdir /usr/local/bin

ENTRYPOINT ["wpdir"]