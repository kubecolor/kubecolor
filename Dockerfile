FROM golang:1.19 as build

WORKDIR /go/src/kubecolor
ADD . /go/src/kubecolor/

RUN go build -o /go/bin/kubecolor cmd/kubecolor/main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/kubecolor /
ENTRYPOINT ["/kubecolor"]