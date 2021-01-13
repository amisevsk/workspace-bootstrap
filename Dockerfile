FROM quay.io/libpod/golang:1.13 AS builder

WORKDIR /workspace-bootstrapper
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY ["main.go", "helper.go", "./"]
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a \
    -o _output/bin/workspace-bootstrap \
    -ldflags="-w -s" \
    main.go helper.go

FROM quay.io/libpod/alpine:latest
WORKDIR /workspace-bootstrapper

RUN apk add --no-cache curl jq git bash && \
    curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

COPY --from=builder /workspace-bootstrapper/_output/bin/workspace-bootstrap ./bootstrap
COPY ["clone-and-sync.sh", "default.devfile.yaml", "./"]

ENV "DEFAULT_DEVFILE" "/workspace-bootstrapper/default.devfile.yaml"

ENTRYPOINT ["/workspace-bootstrapper/clone-and-sync.sh"]
CMD /workspace-bootstrapper/bootstrap
