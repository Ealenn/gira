###################################################
#                     BUILD                       #
###################################################
FROM golang:1.24 AS build
ARG VERSION

WORKDIR /app
COPY . .

RUN echo $VERSION > /app/internal/configuration/version
RUN CGO_ENABLED=0 go build -ldflags="-s -w" /app/cmd/...

###################################################
#                     FINAL                       #
###################################################
FROM alpine:3 AS final

# Install GIT
RUN apk add git

# Copy gira
COPY --from=build /app/gira /usr/bin/gira

WORKDIR /app
ENTRYPOINT [ "gira" ]
