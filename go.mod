module github.com/apache/dubbo-go-pixiu

go 1.18

require (
	cloud.google.com/go/compute v1.6.0
	cloud.google.com/go/security v1.3.0
	contrib.go.opencensus.io/exporter/prometheus v0.4.1
	dubbo.apache.org/dubbo-go/v3 v3.0.2-0.20220519062747-f6405fa79d5c
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20220408101031-f1761e18c0c6
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/MicahParks/keyfunc v1.0.0
	github.com/Shopify/sarama v1.19.0
	github.com/alibaba/sentinel-golang v1.0.4
	github.com/apache/dubbo-getty v1.4.8
	github.com/apache/dubbo-go-hessian2 v1.11.1
	github.com/cch123/supermonkey v1.0.1
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/cncf/xds/go v0.0.0-20220330162227-eded343319d0
	github.com/coreos/go-oidc/v3 v3.1.0
	github.com/creasty/defaults v1.5.2
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v20.10.17+incompatible
	github.com/dubbo-go-pixiu/pixiu-api v0.1.6-0.20220427143451-c0a68bf5b29a
	github.com/dubbogo/go-zookeeper v1.0.4-0.20211212162352-f9d2183d89d5
	github.com/dubbogo/gost v1.11.25
	github.com/dubbogo/grpc-go v1.42.9
	github.com/dubbogo/triple v1.1.8
	github.com/envoyproxy/go-control-plane v0.10.2-0.20220428052930-ec95b9f870a8
	github.com/evanphx/json-patch/v5 v5.6.0
	github.com/fatih/color v1.13.0
	github.com/felixge/fgprof v0.9.2
	github.com/florianl/go-nflog/v2 v2.0.1
	github.com/fsnotify/fsnotify v1.5.4
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.0.1
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/gogo/protobuf v1.3.2
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/cel-go v0.11.5
	github.com/google/go-cmp v0.5.8
	github.com/google/go-containerregistry v0.8.0
	github.com/google/gofuzz v1.2.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.4.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/jhump/protoreflect v1.9.0
	github.com/kr/pretty v0.3.0
	github.com/kylelemons/godebug v1.1.0
	github.com/lestrrat-go/jwx v1.2.23
	github.com/lucas-clemente/quic-go v0.27.0
	github.com/mattn/go-isatty v0.0.14
	github.com/mercari/grpc-http-proxy v0.1.2
	github.com/miekg/dns v1.1.48
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/moby/buildkit v0.10.1
	github.com/nacos-group/nacos-sdk-go v1.1.1
	github.com/onsi/gomega v1.19.0
	github.com/openshift/api v0.0.0-20200713203337-b2494ecb17dd
	github.com/opentrx/seata-golang/v2 v2.0.5
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.33.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/vishvananda/netlink v1.2.0-beta
	github.com/yl2chen/cidranger v1.0.2
	go.etcd.io/etcd/api/v3 v3.6.0-alpha.0
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.6.3
	go.opentelemetry.io/otel/exporters/prometheus v0.21.0
	go.opentelemetry.io/otel/metric v0.30.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/sdk/export/metric v0.21.0
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.opentelemetry.io/proto/otlp v0.15.0
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	golang.org/x/time v0.0.0-20220411224347-583f2d630306
	gomodules.xyz/jsonpatch/v3 v3.0.1
	google.golang.org/api v0.74.0
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.8.2
	istio.io/api v0.0.0-20220728184806-7837c4e62d82
	istio.io/client-go v1.14.3-0.20220728185607-9117649b5a7f
	istio.io/pkg v0.0.0-20220728185106-cbb9bd2a0124
	k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/cli-runtime v0.23.5
	k8s.io/client-go v0.23.5
	k8s.io/klog/v2 v2.60.1
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65
	k8s.io/kubectl v0.23.5
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/controller-runtime v0.11.2
	sigs.k8s.io/gateway-api v0.4.1-0.20220419214231-03f50b47814e
	sigs.k8s.io/mcs-api v0.1.0
	sigs.k8s.io/yaml v1.3.0
	vimagination.zapto.org/byteio v0.0.0-20200222190125-d27cba0f0b10
)

require (
	cloud.google.com/go v0.101.0 // indirect
	cloud.google.com/go/logging v1.5.0-jsonlog-preview // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.12.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/google/pprof v0.0.0-20220729232143-a41b82acbcb1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/klauspost/compress v1.15.8 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.2 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/prometheus/prom2json v1.3.1 // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	go.etcd.io/etcd/server/v3 v3.6.0-alpha.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.starlark.net v0.0.0-20211013185944-b0039bd2cfe3 // indirect
	golang.org/x/tools v0.1.11 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/kustomize/api v0.11.4 // indirect
	vimagination.zapto.org/memio v0.0.0-20200222190306-588ebc67b97d // indirect
)

replace go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v0.21.0

replace go.opentelemetry.io/otel/internal/metric => go.opentelemetry.io/otel/internal/metric v0.21.0

// Client-go does not handle different versions of mergo due to some breaking changes - use the matching version
replace github.com/imdario/mergo v0.3.12 => github.com/imdario/mergo v0.3.5

replace istio.io/api => istio.io/api v0.0.0-20220728184806-7837c4e62d82

replace istio.io/pkg => istio.io/pkg v0.0.0-20220728185106-cbb9bd2a0124
