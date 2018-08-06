# build stage
FROM golang:alpine AS build-env
RUN apk add --update \
    gcc \
    libc-dev \
    ktoblzcheck-dev \
  && rm -rf /var/cache/apk/*
ADD . /src
RUN cd /src && go build -o ktoblzcheck

# app stage
FROM alpine
RUN apk add --update \
    ktoblzcheck \
  && rm -rf /var/cache/apk/*
COPY --from=build-env /src/ktoblzcheck /app/
CMD /app/ktoblzcheck
