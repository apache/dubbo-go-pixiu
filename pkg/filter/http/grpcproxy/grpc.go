package grpcproxy

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"sync"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPGrpcProxyFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// AccessPlugin is http filter plugin.
	Plugin struct {
	}
	// AccessFilter is http filter instance
	Filter struct {
		cfg *Config
	}
	// Config describe the config of AccessFilter
	Config struct{
		Proto string     `yaml:"proto_descriptor" json:"proto_descriptor"`
		rules []*Rule	`yaml:"rules" json:"rules"`
	}

	Rule struct {
		Selector string     `yaml:"selector" json:"selector"`
		Match Match `yaml:"match" json:"match"`
	}

	Match struct {
		method string `yaml:"method" json:"method"`
	}


	fileSource struct {
		files  map[string]*desc.FileDescriptor
		er     *dynamic.ExtensionRegistry
		erInit sync.Once
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (af *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(af.Handle)
	return nil
}

// Handle use the default http to grpc transcoding strategy https://cloud.google.com/endpoints/docs/grpc/transcoding
func (af *Filter) Handle(c *http.HttpContext) {

}

func (af *Filter) Config() interface{} {
	return af.cfg
}

func (af *Filter) Apply() error {

}

func (af *Filter) initFromFileDescriptor(importPaths []string, fileNames ...string) error{
	fileNames, err := protoparse.ResolveFilenames(importPaths, fileNames...)
	if err != nil {
		return err
	}
	p := protoparse.Parser{
		ImportPaths:           importPaths,
		InferImportPaths:      len(importPaths) == 0,
		IncludeSourceCodeInfo: true,
	}
	fds, err := p.ParseFiles(fileNames...)
	if err != nil {
		return fmt.Errorf("could not parse given files: %v", err)
	}
	return nil
}

func DescriptorSourceFromFileDescriptors(files ...*desc.FileDescriptor) (*fileSource, error) {
	fds := map[string]*desc.FileDescriptor{}
	for _, fd := range files {
		if err := addFile(fd, fds); err != nil {
			return nil, err
		}
	}
	return &fileSource{files: fds}, nil
}


func addFile(fd *desc.FileDescriptor, fds map[string]*desc.FileDescriptor) error {
	name := fd.GetName()
	if existing, ok := fds[name]; ok {
		// already added this file
		if existing != fd {
			// doh! duplicate files provided
			return fmt.Errorf("given files include multiple copies of %q", name)
		}
		return nil
	}
	fds[name] = fd
	for _, dep := range fd.GetDependencies() {
		if err := addFile(dep, fds); err != nil {
			return err
		}
	}
	return nil
}

func (fs *fileSource) ListServices() ([]string, error) {
	set := map[string]bool{}
	for _, fd := range fs.files {
		for _, svc := range fd.GetServices() {
			set[svc.GetFullyQualifiedName()] = true
		}
	}
	sl := make([]string, 0, len(set))
	for svc := range set {
		sl = append(sl, svc)
	}
	return sl, nil
}

func (fs *fileSource) GetAllFiles() ([]*desc.FileDescriptor, error) {
	files := make([]*desc.FileDescriptor, len(fs.files))
	i := 0
	for _, fd := range fs.files {
		files[i] = fd
		i++
	}
	return files, nil
}

func (fs *fileSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	for _, fd := range fs.files {
		if dsc := fd.FindSymbol(fullyQualifiedName); dsc != nil {
			return dsc, nil
		}
	}
	return nil, fmt.Errorf("could not found symbol %v", fullyQualifiedName)
}

func (fs *fileSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	fs.erInit.Do(func() {
		fs.er = &dynamic.ExtensionRegistry{}
		for _, fd := range fs.files {
			fs.er.AddExtensionsFromFile(fd)
		}
	})
	return fs.er.AllExtensionsForType(typeName), nil
}



