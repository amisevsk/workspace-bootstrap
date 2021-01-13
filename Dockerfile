FROM quay.io/libpod/alpine:latest

RUN apk add --no-cache curl jq git bash python3 py-pip && \
    pip install yq && \
    curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

COPY clone-and-sync.sh /clone-and-sync.sh
CMD /clone-and-sync.sh
