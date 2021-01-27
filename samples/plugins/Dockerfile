FROM golang:1.14

ENV HOME /go/src
ENV GOPROXY https://goproxy.cn

WORKDIR ${HOME}

COPY . ${HOME}

CMD go version\
&& go build -o out -buildmode=plugin ${HOME}/plugin/*.go\
&& echo "build sucessful!\n"\
