FROM alpine:3.9 as base

FROM golang:1.12 as builder

ENV go_workspace=/go
ENV GOPATH $go_workspace
ENV PROJECT_DIR=$GOPATH/src/github.com/sagansystems/terraform-provider-snowflake
WORKDIR ${PROJECT_DIR}

COPY snowflake snowflake
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN gofmt -l `find . -name '*.go'`
RUN GO111MODULE=on go build

# back to the base-ics, we only need the binary
FROM base

ENV BIN_DIR=/go/src/github.com/sagansystems/terraform-provider-snowflake
WORKDIR /usr/local/bin

COPY --from=builder /go/src/github.com/sagansystems/terraform-provider-snowflake/terraform-provider-snowflake /usr/local/bin/
