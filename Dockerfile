FROM golang:1.21 as build

WORKDIR /go/src/kubecolor
ADD . /go/src/kubecolor/

RUN go build -o /go/bin/kubecolor cmd/kubecolor/main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/kubecolor /
ENTRYPOINT ["/kubecolor"]

LABEL org.opencontainers.image.source=https://github.com/kubecolor/kubecolor
LABEL org.opencontainers.image.description="Colorize your kubectl output"
LABEL org.opencontainers.image.licenses=MIT
