# build stage
FROM golang:1.10.3 AS build-env
RUN apt-get update -y && apt-get install wget -y

RUN mkdir ktoblzcheck
WORKDIR ktoblzcheck
RUN curl -L https://sourceforge.net/projects/ktoblzcheck/files/ktoblzcheck-1.49.tar.gz/download | tar zxv
RUN cd ktoblzcheck-1.49 && ./configure && make && make install
ADD . /src
RUN cd /src && go build -o ktoblzcheck

# app stage
FROM alpine
COPY --from=build-env /src/ktoblzcheck /app/
CMD /app/ktoblzcheck
