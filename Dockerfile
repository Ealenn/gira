###################################################
#                     BUILD                       #
###################################################
FROM golang:1.24 AS build
ARG VERSION

WORKDIR /app
COPY . .

RUN echo $VERSION > /app/configuration/version.txt
RUN CGO_ENABLED=0 go build -o gira -ldflags="-s -w"

###################################################
#                     FINAL                       #
###################################################
FROM alpine:3

# Install GIT
RUN apk add git

# Copy gira
COPY --from=build /app/gira /usr/bin/gira

WORKDIR /app
ENTRYPOINT [ "gira" ]
