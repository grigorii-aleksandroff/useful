FROM harbor-adm-dc1.itechpsp.com/dockerhub/golang:1.25.1-alpine AS go-builder

ARG BOX_CI_USER
ARG BOX_CI_TOKEN
ARG ITECH_BOX_CI_USER
ARG ITECH_BOX_CI_TOKEN

RUN apk update && apk add --no-cache git curl bash

RUN echo "machine git.itechpsp.com login ${ITECH_BOX_CI_USER} password ${ITECH_BOX_CI_TOKEN}" > ~/.netrc
RUN echo "machine git.bububla.com login ${BOX_CI_USER} password ${BOX_CI_TOKEN}" >> ~/.netrc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN go build -ldflags="-X main.commit=$(git rev-parse HEAD) -X main.tag=$(git describe --tags --always)" -o bin/main main.go

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

FROM harbor-adm-dc1.itechpsp.com/dockerhub/alpine:latest

WORKDIR /app

COPY --from=go-builder /app/bin /app/bin
COPY --from=go-builder /go/bin/grpcurl /usr/local/bin/grpcurl

ENV PATH="/app/bin:${PATH}"

CMD ["main"]