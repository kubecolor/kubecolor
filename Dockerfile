FROM docker.io/library/golang:1.24.7 AS build

WORKDIR /go/src/kubecolor
COPY go.mod go.sum .
RUN go mod download

COPY . .
ARG VERSION
RUN CGO_ENABLED=0 go install -ldflags="-X main.Version=${VERSION}" .

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /go/bin/kubecolor /usr/local/bin/
COPY --from=bitnami/kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/
ENTRYPOINT ["kubecolor"]

LABEL org.opencontainers.image.source=https://github.com/kubecolor/kubecolor
LABEL org.opencontainers.image.description="Colorize your kubectl output"
LABEL org.opencontainers.image.licenses=MIT
