package grpcproxy

import (
	"fmt"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

type fileSource struct {
	files  map[string]*desc.FileDescriptor
	er     *dynamic.ExtensionRegistry
	erInit sync.Once
}

// initFromFileDescriptor
// protoparse.ResolveFilenames(importPaths, fileNames...)
// rel: https://pkg.go.dev/github.com/jhump/protoreflect/desc/protoparse?utm_source=godoc#ResolveFilenames
func (af *Filter) initFromFileDescriptor(importPaths []string, fileNames ...string) error {
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

	for _, fd := range fds {
		name := fd.GetName()
		fs.files[name] = fd
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
