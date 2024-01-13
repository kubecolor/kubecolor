FROM golang:1.21.6 as build

WORKDIR /go/src/kubecolor
COPY go.mod go.sum .
RUN go mod download

COPY . .
ARG VERSION
RUN go install -ldflags="-X main.Version=${VERSION}" .

FROM gcr.io/distroless/base:nonroot
COPY --from=build /go/bin/kubecolor /usr/local/bin/
COPY --from=bitnami/kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/
ENTRYPOINT ["kubecolor"]

LABEL org.opencontainers.image.source=https://github.com/kubecolor/kubecolor
LABEL org.opencontainers.image.description="Colorize your kubectl output"
LABEL org.opencontainers.image.licenses=MIT
