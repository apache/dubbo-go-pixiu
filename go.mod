module github.com/apache/dubbo-go-pixiu

go 1.18

// https://github.com/containerd/containerd/issues/5781
exclude k8s.io/kubernetes v1.13.0

// Client-go does not handle different versions of mergo due to some breaking changes - use the matching version
replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.5

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
	github.com/apache/dubbo-go-hessian2 v1.11.3
	github.com/cch123/supermonkey v1.0.1
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/cncf/xds/go v0.0.0-20220330162227-eded343319d0
	github.com/coreos/go-oidc/v3 v3.1.0
	github.com/creasty/defaults v1.5.2
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/cli v20.10.17+incompatible
	github.com/dubbo-go-pixiu/pixiu-api v0.1.6-0.20220612115254-d9a176b25b99
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
	github.com/imdario/mergo v0.3.12
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
	github.com/prometheus/client_golang v1.13.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.37.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.2
	github.com/vishvananda/netlink v1.2.0-beta
	github.com/yl2chen/cidranger v1.0.2
	go.etcd.io/etcd/api/v3 v3.5.5
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.10.0
	go.opentelemetry.io/otel/exporters/prometheus v0.32.1
	go.opentelemetry.io/otel/metric v0.32.1
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/sdk/metric v0.32.1
	go.opentelemetry.io/otel/trace v1.10.0
	go.opentelemetry.io/proto/otlp v0.19.0
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a
	golang.org/x/time v0.0.0-20220411224347-583f2d630306
	gomodules.xyz/jsonpatch/v3 v3.0.1
	google.golang.org/api v0.74.0
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.1
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	helm.sh/helm/v3 v3.9.4
	istio.io/api v0.0.0-20220906212632-fa18ce1fc2da
	istio.io/client-go v1.14.4-0.20220906213432-736be6f83263
	istio.io/pkg v0.0.0-20220906212832-f98e656e3df0
	k8s.io/api v0.24.2
	k8s.io/apiextensions-apiserver v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/cli-runtime v0.24.2
	k8s.io/client-go v0.24.2
	k8s.io/klog/v2 v2.60.1
	k8s.io/kube-openapi v0.0.0-20220627174259-011e075b9cb8
	k8s.io/kubectl v0.24.2
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	mosn.io/proxy-wasm-go-host v0.1.0
	sigs.k8s.io/controller-runtime v0.11.2
	sigs.k8s.io/gateway-api v0.4.1-0.20220419214231-03f50b47814e
	sigs.k8s.io/mcs-api v0.1.0
	sigs.k8s.io/yaml v1.3.0
	vimagination.zapto.org/byteio v0.0.0-20200222190125-d27cba0f0b10
)

require (
	cloud.google.com/go v0.101.0 // indirect
	cloud.google.com/go/logging v1.5.0-jsonlog-preview // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20170808103936-bb23615498cd // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/RoaringBitmap/roaring v0.7.1 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/Workiva/go-datastructures v1.0.52 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.18 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220418222510-f25a4f6275ed // indirect
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chai2010/gettext-go v0.0.0-20160711120539-c6fed771bfd5 // indirect
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.11.3 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fvbommel/sortorder v1.0.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/pprof v0.0.0-20220729232143-a41b82acbcb1 // indirect
	github.com/googleapis/gax-go/v2 v2.3.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v0.0.0-20200817173448-b6b71def0850 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible // indirect
	github.com/klauspost/compress v1.15.8 // indirect
	github.com/knadh/koanf v1.4.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/marten-seemann/qpack v0.2.1 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.5 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mdlayher/netlink v1.4.1 // indirect
	github.com/mdlayher/socket v0.0.0-20210307095302-262dc9984e00 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/prom2json v1.3.1 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect
	github.com/shirou/gopsutil v3.21.3+incompatible // indirect
	github.com/shirou/gopsutil/v3 v3.21.10 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/wasmerio/wasmer-go v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.5 // indirect
	go.etcd.io/etcd/client/v3 v3.5.5 // indirect
	go.etcd.io/etcd/server/v3 v3.5.5 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.starlark.net v0.0.0-20211013185944-b0039bd2cfe3 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/arch v0.0.0-20200826200359-b19915210f00 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	gomodules.xyz/orderedmap v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	k8s.io/component-base v0.24.2 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/kustomize/api v0.11.4 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	vimagination.zapto.org/memio v1.0.0 // indirect
)

replace google.golang.org/protobuf v1.28.1 => google.golang.org/protobuf v1.28.0
