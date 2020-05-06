FROM golang

WORKDIR /go/src/github.com/thuc201995/srv_convert_file

COPY . .

RUN go get  -u ./...
RUN go install -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

RUN apt-get update -y &&\
    apt-get install software-properties-common -y &&\
    add-apt-repository -r ppa:libreoffice/ppa -y &&\
    apt-get update -q &&\
    apt install libreoffice -y

ENTRYPOINT  watcher -run github.com/thuc201995/srv_convert_file/cmd  -watch github.com/thuc201995/srv_convert_file
