# How run dubbo simple test samples

## Start pixiu

- build pixiu

```bash
cd /dubbogo-go-pixiu/cmd/pixiu

go build
```

- move pixiu to run folder

```bash
mv -b pixiu ../../samples/dubbogo/simple
```

- start one sample

Use `./start.sh [name]` to start pixiu，you have a few choices：

```bash
./start.sh body
./start.sh mix
./start.sh query
./start.sh uri
./start.sh pixiu
```

- interface access

Go to the path of the specified test example，such as `/samples/dubbogo/simple/body`

Executive Command：

```bash
./request.sh
```

