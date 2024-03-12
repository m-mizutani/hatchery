FROM golang:1.22 AS build-go
ENV CGO_ENABLED=0
ARG BUILD_VERSION
COPY . /app
WORKDIR /app
RUN go build -o hatchery -ldflags "-X github.com/m-mizutani/hatchery/pkg/domain/model.AppVersion=${BUILD_VERSION}" .

RUN curl -L -o /pkl https://github.com/apple/pkl/releases/download/0.25.2/pkl-linux-amd64
RUN chmod +x /pkl

FROM --platform=linux/x86_64 ubuntu:20.04
COPY --from=build-go /app/hatchery /hatchery
COPY --from=build-go /pkl /usr/local/bin/pkl

ENTRYPOINT ["/hatchery"]
