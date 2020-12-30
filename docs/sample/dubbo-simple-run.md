# How run dubbo simple test samples

## Start proxy

- build proxy

```bash
cd /dubbogo-go-proxy/cmd/proxy

go build
```

- move proxy to run folder

```bash
mv -b proxy ../../samples/dubbogo/simple
```

- start one sample

Use `./start.sh [name]` to start proxy，you have a few choices：

```bash
./start.sh body
./start.sh mix
./start.sh query
./start.sh uri
./start.sh proxy
```

- interface access

Go to the path of the specified test example，such as `/samples/dubbogo/simple/body`

Executive Command：

```bash
./request.sh
```

