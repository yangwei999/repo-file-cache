FROM golang:latest as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
COPY . /go/src/github.com/opensourceways/repo-file-cache
RUN cd /go/src/github.com/opensourceways/repo-file-cache && CGO_ENABLED=1 go build -v -o ./repo-file-cache main.go

# copy binary config and utils
FROM golang:latest
RUN  mkdir -p /opt/app/
RUN  mkdir -p /opt/app/controllers
COPY ./go.mod /opt/app/
COPY ./go.sum /opt/app/
COPY ./conf /opt/app/conf
# overwrite config yaml
COPY  --from=BUILDER /go/src/github.com/opensourceways/repo-file-cache/repo-file-cache /opt/app

WORKDIR /opt/app/
ENTRYPOINT ["/opt/app/repo-file-cache"]
