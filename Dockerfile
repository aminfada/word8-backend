FROM golang:1.20 AS builder
COPY . /build/
WORKDIR /build
RUN go mod tidy
RUN go mod vendor
RUN go build -o main /build/


FROM ubuntu:20.04
# Fix certificate issues
RUN apt-get update -y && \
    apt-get install ca-certificates-java -y && \
    apt-get clean -y && \
    update-ca-certificates -f -y;

COPY --from=builder /build/main /app/

ADD ./docs/db /app/docs/db
ADD ./docs/words /app/docs/words

ENTRYPOINT [ "/app/main" ]
