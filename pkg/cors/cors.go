package cors

import "net/http"

// https://github.com/rs/cors/blob/master/cors.go

type Options struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowOriginFunc is a custom function to validate the origin. It take the
	// origin as argument and returns true if allowed or false otherwise. If
	// this option is set, the content of `AllowedOrigins` is ignored.
	AllowOriginFunc func(origin string) bool
	// AllowOriginRequestFunc is a custom function to validate the origin. It
	// takes the HTTP Request object and the origin as argument and returns true
	// if allowed or false otherwise. If headers are used take the decision,
	// consider using AllowOriginVaryRequestFunc instead. If this option is set,
	// the contents of `AllowedOrigins`, `AllowOriginFunc` are ignored.
	//
	// Deprecated: use `AllowOriginVaryRequestFunc` instead.
	AllowOriginRequestFunc func(r *http.Request, origin string) bool
	// AllowOriginVaryRequestFunc is a custom function to validate the origin.
	// It takes the HTTP Request object and the origin as argument and returns
	// true if allowed or false otherwise with a list of headers used to take
	// that decision if any so they can be added to the Vary header. If this
	// option is set, the contents of `AllowedOrigins`, `AllowOriginFunc` and
	// `AllowOriginRequestFunc` are ignored.
	AllowOriginVaryRequestFunc func(r *http.Request, origin string) (bool, []string)
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [].
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached. Default value is 0, which stands for no
	// Access-Control-Max-Age header to be sent back, resulting in browsers
	// using their default value (5s by spec). If you need to force a 0 max-age,
	// set `MaxAge` to a negative value (ie: -1).
	MaxAge int
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// AllowPrivateNetwork indicates whether to accept cross-origin requests over a
	// private network.
	AllowPrivateNetwork bool
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
	// Provides a status code to use for successful OPTIONS requests.
	// Default value is http.StatusNoContent (204).
	OptionsSuccessStatus int
	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool
}

// Cors http handler
type Cors struct {
	// Normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	// allowedWOrigins []wildcard
	// Optional origin validator function
	allowOriginFunc func(r *http.Request, origin string) (bool, []string)
	// Normalized list of allowed headers
	// Note: the Fetch standard guarantees that CORS-unsafe request-header names
	// (i.e. the values listed in the Access-Control-Request-Headers header)
	// are unique and sorted;
	// see https://fetch.spec.whatwg.org/#cors-unsafe-request-header-names.
	// allowedHeaders internal.SortedSet
	// Normalized list of allowed methods
	allowedMethods []string
	// Pre-computed normalized list of exposed headers
	exposedHeaders []string
	// Pre-computed maxAge header value
	maxAge []string
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	// Status code to use for successful OPTIONS requests
	optionsSuccessStatus int
	allowCredentials     bool
	allowPrivateNetwork  bool
	optionPassthrough    bool
	preflightVary        []string
}
