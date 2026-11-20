# Start

How to start the dubbo-go-pixiu

#### 1 cd samples dir

```
cd samples/dubbo/simple
```

we can use start.sh to run samples quickly. for more info, execute command as below for more help

```
./start.sh [action] [project]
./start.sh help
```

we run body samples below step

#### 2 prepare config file and docker

prepare command will prepare dubbo-server and pixiu config file and start docker container needed

```
./start.sh prepare body
```

if prepare config file manually, notice:

- modify $PROJECT_DIR in conf.yaml to absolute path in your compute

#### 3 start dubbo or http server

```
./start.sh startServer body
```

#### 4 start pixiu

```
./start.sh startPixiu body
```

if run pixiu manually, use command as below

```
 go run cmd/pixiu/*.go gateway start -c /[absolute-path]/dubbo-go-pixiu/samples/dubbo/simple/body/pixiu/conf.yaml
```

#### 5. Try a request

use curl to try or use unit test

```bash
curl -X POST 'localhost:8881/api/v1/test-dubbo/user' -d '{"id":"0003","code":3,"name":"dubbogo","age":99}' --header 'Content-Type: application/json' 
./start.sh startTest body
```

#### 6. Clean

```
./start.sh clean body
```

## Graceful shutdown and exit status

Enable signal-driven graceful shutdown with a positive timeout in the gateway
bootstrap configuration:

```yaml
static_resources:
  shutdown_config:
    timeout: "60s"
```

On SIGTERM or SIGINT, Pixiu starts all listener shutdowns concurrently and waits
for their return values, bounded by this overall timeout. There is no shorter,
fixed per-listener deadline.

| Shutdown result | Process exit code |
| --- | --- |
| Every listener returns without error within the timeout | `0` |
| Any listener returns an error | `1` |
| The overall timeout expires | `1` |

**Upgrade note:** shutdown errors and timeouts previously exited with code `0`.
The nonzero status is intentional: receiving SIGTERM does not make an incomplete
shutdown successful. Kubernetes or systemd may record a failed termination, and
restart or alerting rules that depend on exit status may react to it. Review
those rules before upgrading.

If shutdown exits with code `1`, inspect the gateway logs before the termination
for `Shutdown Error:`, `Shutdown completed with`, and `Shutdown timed out after`.
Check listener cleanup failures and whether active requests can finish within
the configured timeout. Allow more time than Pixiu's timeout in Kubernetes
`terminationGracePeriodSeconds` or systemd `TimeoutStopSec`, including time spent
in any termination hooks. An external forced kill can prevent Pixiu from
finishing cleanup or reporting its own exit status.

A timeout of `0s` (the default when `shutdown_config` is omitted), or a negative
timeout, **disables** Pixiu's graceful shutdown signal handler. OS signal behavior
then applies; the manager does not invoke listener `ShutDown()` or `Close()` on
termination, and the exit-code table above does not apply. Zero does not mean
"close immediately and run cleanup". Use a positive timeout when cleanup is
required.
