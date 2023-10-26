// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

const (
	Ionos_tokenScopes = "ionos_token.Scopes"
)

// DBaaSQuota defines model for DBaaSQuota.
type DBaaSQuota struct {
	CPU              *int64 `json:"CPU,omitempty"`
	Memory           *int64 `json:"Memory,omitempty"`
	MongoClusters    *int64 `json:"MongoClusters,omitempty"`
	PostgresClusters *int64 `json:"PostgresClusters,omitempty"`
	Storage          *int64 `json:"Storage,omitempty"`
}

// DNSQuota defines model for DNSQuota.
type DNSQuota struct {
	Records        *int64 `json:"Records,omitempty"`
	SecondaryZones *int64 `json:"SecondaryZones,omitempty"`
	Zones          *int64 `json:"Zones,omitempty"`
}

// Quotas defines model for Quotas.
type Quotas struct {
	DBaaS *struct {
		Limits *DBaaSQuota `json:"Limits,omitempty"`
		Usage  *DBaaSQuota `json:"Usage,omitempty"`
	} `json:"DBaaS,omitempty"`
	DNS *struct {
		Limits *DNSQuota `json:"Limits,omitempty"`
		Usage  *DNSQuota `json:"Usage,omitempty"`
	} `json:"DNS,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	GetHealth(w http.ResponseWriter, r *http.Request)

	// (GET /quotas)
	GetQuotas(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /health)
func (_ Unimplemented) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /quotas)
func (_ Unimplemented) GetQuotas(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetQuotas operation middleware
func (siw *ServerInterfaceWrapper) GetQuotas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Ionos_tokenScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetQuotas(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.GetHealth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/quotas", wrapper.GetQuotas)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xVTXPTMBD9K5qFo4lTUjj4VloKoZAEQocZMjmoysZWa2tdaZ0h0/F/ZySlDW3S4jKc",
	"rI/9ePu0b30DiqqaDBp2kN2ARVeTcRg2760l6xeKDKNhv2T8xWldSm38zqkCKxnO1zVCBo6tNjm0bZvA",
	"Ap2yumZNBjI4mgyFkmUpMERtE/jaEEv3IL6s61Ir6Z3SS0cPsry0uIQMXqRb1Gm8dekm3L7USjVVU0rG",
	"hbgOVmJJVngwDu1KK/RuySZPQHTyTsppiOh3taUaLetIy/Hk3H+WZCvJkIE2/PYQklsKtGHMMVT4BSuy",
	"667GZHI6LhvHaF1Hnwk5zi26Z7pNmazMsZN1e3dEF5eo2PufjB7j5hsqsovOOFCRWUi7/kkGuzp1t90H",
	"fdt094GHB989/qwrHZXxVOv90S1tAuduw21Xl0co/mc0o+diGT2OZO+JQ9VYzeupDxChaTLkXjFdYZDs",
	"BUqL9vT2hT79+A4bfflQ8Xb7YgVzHYWrzZK8/30BnyKrAp2QeW4x39GxyPUKjZBKUWNYLC1VggsUw/Fo",
	"PBXHJTULcTQZ+nyaS58wdoGYRvVDAiu0Lubq9w56fU8d1WhkrSGDQa/fG0ACteQiFJsWKEsu/DJH3sX7",
	"MVwLVaC6ghDJhok2XEAGH5DjPST3x+3rfr/jsEXTVJDNYHwG8+Tvo3d85ut5E8Pva4Q7GGkc+W0Ikl7f",
	"aWVvlf/tVXb42Yh0Pz9PF3D7F0jgsH/Qud5nsrNVAGSzB70/m7fzNprYVRjJsxtobLlp8yxNS1KyLMhx",
	"NhgMBtDO298BAAD//2Eoh3OGBwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
