FROM gcr.io/istio-testing/build-tools:latest as build-env
RUN go get github.com/go-delve/delve/cmd/dlv

FROM docker.io/istio/pilot:1.14.3

COPY main /usr/local/bin/pilot-discovery
COPY --from=build-env /gobin/dlv /usr/local/bin/dlv
EXPOSE 40000
ENTRYPOINT ["/usr/local/bin/dlv", "--continue", "--accept-multiclient", "--listen=:40000", "--check-go-version=false", "--headless=true", "--api-version=2", "--log=true", "--log-output=debugger,debuglineerr,gdbwire,lldbout,rpc", "exec", "/usr/local/bin/pilot-discovery", "--"]