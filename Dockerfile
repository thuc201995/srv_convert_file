FROM golang

RUN mkdir -p /go/src/github.com/thuc201995

ADD . /go/src/github.com/thuc201995

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run github.com/thuc201995/srv_convert_file/cmd  -watch github.com/thuc201995/srv_convert_file
