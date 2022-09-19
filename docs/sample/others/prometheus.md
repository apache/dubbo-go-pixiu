# Metric Filter Quick Start

## Start PrometheusPushGateway [Docker environment]:

##### Use Docker to install and start：

 Directly obtain the latest version of the official image `prom/pushgateway:latest` The startup command is as follows:

```shell
$ docker pull prom/pushgateway
```

```shell
$ docker run -d -p 9091:9091 prom/pushgateway
```

Use the command `./pushgateway` command to start the service. At this time, the browser can access the UI page by accessing `http://<ip>:9091`, but there is no data display on the default Metrics, that is because we have not yet sent to PushGateway Push any data. 

However, the PushGateway service itself comes with some metrics, which can be obtained by visiting the `http://<ip>:9091/metrics` address. You can see that it contains some monitoring metrics related to go and process.

##### Check if the configuration is successful：

PushGateway provides a standard API interface and allows users to add data. The default URL address is: `http://<ip>:9091/metrics/job/<JOBNAME>{/<LABEL_NAME>/<LABEL_VALUE>}`, where `<JOBNAME>` It is a required item, which is the value of the job label. It can be followed by any number of label pairs. Generally, we will add an instance`/<INSTANCE_NAME>` instance name label to facilitate the distinction of each indicator.  Next, you can push a simple indicator data to PushGateway for testing.

```shell
$ echo "test_metric 123456" | curl --data-binary @- http://<ip>:9091/metrics/job/test_job
```

After the execution is complete, refresh the PushGateway UI page to verify that you can see the test_metric indicator data you just added.

It can also be tested in the following way:
```shell
$ cat <<EOF | curl --data-binary @- http://<ip>:9091/metrics/job/test_job/instance/test_instance
# TYPE test_metrics counter
test_metrics{label="app1",name="demo"} 100.00
# TYPE another_test_metrics gauge
# HELP another_test_metrics Just an example.
another_test_metrics 123.45
EOF
```

## Start Pixiu:

Examples of official references is in `https://github.com/dubbo-go-pixiu/samples`

Add the following configuration file to the `samples/http/simple/pixiu/conf.yaml`

```yaml
static_resources:
  listeners:
    - name: "net/http"
      protocol_type: "HTTP"
      address:
        socket_address:
          address: "0.0.0.0"
          port: 8888
      filter_chains:
        filters:
          - name: dgp.filter.httpconnectionmanager
            config:
              route_config:
                routes:
                  - match:
                      prefix: /user
                    route:
                      cluster: user
                      cluster_not_found_response_code: 505
                http_filters:
                - name: dgp.filter.http.prometheusmetric
                  metric_collect_rules:
                    enable: true
                    metric_path: "/metrics"
                    push_gateway_url: "http://127.0.0.1:9091"
                    push_interval_seconds: 3
                    push_job_name: "prometheus"
      config:
        idle_timeout: 5s
        read_timeout: 5s
        write_timeout: 5s
  clusters:
    - name: "user"
      lb_policy: "lb"
      endpoints:
        - id: 1
          socket_address:
            address: 127.0.0.1
            port: 1314
      health_checks:
        - protocol: "tcp"
          timeout: 1s
          interval: 2s
          healthy_threshold: 4
          unhealthy_threshold: 4
  shutdown_config:
    timeout: "60s"
    step_timeout: "10s"
    reject_policy: "immediacy"
```

Then execute the following command .

```shell
go run cmd/pixiu/*.go gateway start -c samples/http/simplep/pixiu/conf.yaml
```
Then you can also query the collected indicator data on the PushGateway UI page .
