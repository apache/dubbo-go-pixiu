### xds

xds implement demo how to use xds server.

### to run sample:

1. start xds server
```
dubbo-go-pixiu/xds/xdsserver > go run .
```

2. run pixiu 
```
dubbo-go-pixiu > pixiu gateway start -c ./samples/xds/pixiu/conf.yaml -g test/configs/log.yml
```

3. check result
```shell
curl -v  'localhost:8889/get'

## will get below 
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.64.1",
    "X-Amzn-Trace-Id": "Root=1-61ba16a5-3ea1961217b2ffa7124ea2c2"
  },
  "origin": "223.104.41.209",
  "url": "http://httpbin.org/get"
}
```