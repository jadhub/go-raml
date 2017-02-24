package golang

import (
	"path/filepath"
	"strings"

	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/codegen/resource"
)

// ClientService represents a root endpoint of an API
type ClientService struct {
	rootEndpoint string
	PackageName  string
	Methods      []resource.MethodInterface
}

// Name returns it's struct name
func (cs ClientService) Name() string {
	name := commons.NormalizeName(strings.Title(cs.rootEndpoint[1:]) + "Service")

	return name
}

// EndpointName returns root endpoint name
func (cs ClientService) EndpointName() string {
	result := commons.NormalizeName(strings.Title(cs.rootEndpoint[1:]))

	return result
}

// FilenameNoExt return filename without extension
func (cs ClientService) FilenameNoExt() string {
	return cs.rootEndpoint[1:] + "_service"
}
func (cs ClientService) filename(dir string) string {
	name := filepath.Join(dir, cs.FilenameNoExt())
	return name + ".go"
}

// LibImportPaths returns all imported lib
func (cs ClientService) LibImportPaths() map[string]struct{} {
	ip := map[string]struct{}{}

	// methods
	for _, v := range cs.Methods {
		gm := v.(clientMethod)
		for lib := range gm.libImported(globRootImportPath) {
			ip[lib] = struct{}{}
		}
	}
	return ip
}
