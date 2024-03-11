FROM golang:1.22 AS build-go
ENV CGO_ENABLED=0
ARG BUILD_VERSION
COPY . /app
WORKDIR /app
RUN go build -o hatchery -ldflags "-X github.com/m-mizutani/hatchery/pkg/domain/model.AppVersion=${BUILD_VERSION}" .

FROM gcr.io/distroless/base:nonroot
USER nonroot
COPY --from=build-go /app/hatchery /hatchery

ENTRYPOINT ["/hatchery"]
