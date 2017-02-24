package golang

import (
	"github.com/Jumpscale/go-raml/codegen/commons"
	"github.com/Jumpscale/go-raml/raml"
)

// generate all body struct from an RAML definition
func generateBodyStructs(apiDef *raml.APIDefinition, dir, packageName string) error {
	// generate
	for _, v := range apiDef.Resources {
		if err := generateStructsFromResourceBody("", dir, packageName, &v); err != nil {
			return err
		}
	}

	return nil
}

// generate all structs from resource's method's request & response body
func generateStructsFromResourceBody(resourcePath, dir, packageName string, r *raml.Resource) error {
	if r == nil {
		return nil
	}

	// build
	var methods = []struct {
		Name   string
		Method *raml.Method
	}{
		{Name: "Get", Method: r.Get},
		{"Post", r.Post},
		{"Head", r.Head},
		{"Put", r.Put},
		{"Delete", r.Delete},
		{"Patch", r.Patch},
		{"Options", r.Options},
	}

	normalizedPath := commons.NormalizeURITitle(resourcePath + r.URI)
	normalizedPath = commons.NormalizeName(normalizedPath)

	for _, v := range methods {
		if err := buildBodyFromMethod(normalizedPath, v.Name, dir, packageName, v.Method); err != nil {
			return err
		}
	}

	// build request/response body of child resources
	for _, v := range r.Nested {
		if err := generateStructsFromResourceBody(resourcePath+r.URI, dir, packageName, v); err != nil {
			return err
		}
	}

	return nil
}

// build request and reponse body of a method.
// in python case, we only need to build it for request body because we only need it for validator
func buildBodyFromMethod(normalizedPath, methodName, dir, packageName string, method *raml.Method) error {
	if method == nil {
		return nil
	}

	//generate struct for request body
	if err := generateStructFromBody(normalizedPath+methodName, dir, packageName, &method.Bodies, true); err != nil {
		return err
	}

	//generate struct for response body
	for _, val := range method.Responses {
		if err := generateStructFromBody(normalizedPath+methodName, dir, packageName, &val.Bodies, false); err != nil {
			return err
		}

	}

	return nil
}

// generate a struct from an RAML request/response body
func generateStructFromBody(structNamePrefix, dir, packageName string, body *raml.Bodies, isGenerateRequest bool) error {
	if !commons.HasJSONBody(body) {
		return nil
	}

	// construct struct from body
	structDef := newStructDefFromBody(body, structNamePrefix, packageName, isGenerateRequest)

	// generate
	return structDef.generate(dir)
}
