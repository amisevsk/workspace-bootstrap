FROM quay.io/libpod/golang:1.13 AS builder

WORKDIR /workspace-sync
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY ["devworkspace-sync/main.go", "./devworkspace-sync/"]
COPY ["devworkspace-sync/funcs", "./devworkspace-sync/funcs"]
COPY ["library", "./library"]
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a \
    -o _output/bin/workspace-sync \
    -ldflags="-w -s" \
    ./devworkspace-sync/main.go

FROM quay.io/libpod/alpine:latest
WORKDIR /workspace-sync

# Useful for debugging
RUN apk add --no-cache curl jq git bash

COPY --from=builder /workspace-sync/_output/bin/workspace-sync /usr/local/bin/workspace-sync

COPY devworkspace-sync/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

USER ${USER_UID}

ENTRYPOINT ["/usr/local/bin/entrypoint"]
CMD ["tail", "-f", "/dev/null"]
