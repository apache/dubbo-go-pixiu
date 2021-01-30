FROM golang:1.14

ENV HOME /go/src
ENV GOPROXY https://goproxy.cn

WORKDIR ${HOME}

COPY . ${HOME}

CMD ["/bin/bash", "${HOME}/start.sh"]