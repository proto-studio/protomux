package mux

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"proto.zip/studio/mux/pkg/muxcontext"
	"proto.zip/studio/validate/pkg/errors"
)

// Error handler function interface for router.Mux implementations that use the
// standard HTTP server method.
type HttpErrorHandler func(error, http.ResponseWriter, *http.Request)

// Implementation of the router.Mux pattern using standard HTTP server method.
type HttpMux struct {
	Mux[http.Handler, HttpErrorHandler]
}

// Implementation of the error interface for HTTP specific errors to
// allow better more semantic API responses.
type HttpError struct {
	// The HTTP status code for the error.
	StatusCode int
}

// Implements standard error response for HttpError.
func (err HttpError) Error() string {
	return http.StatusText(err.StatusCode)
}

// Creates a new HttpError with a specific status.
func NewHttpError(code int) HttpError {
	return HttpError{
		StatusCode: code,
	}
}

// Default HTTP error handler.
//
// If the error is an HttpError it will serve the status text as a string.
// If it is a validation error on path or host it will return 404.
// If it is a validation error on query string or body it will return 400.
// Otherwise it will log the error stack and return a 500 error.
func DefaultErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	switch errCast := err.(type) {
	case HttpError:
		w.WriteHeader(errCast.StatusCode)
		w.Write([]byte(err.Error()))
	case errors.ValidationError:
		w.WriteHeader(400)
		w.Write([]byte(http.StatusText(400)))
	default:
		fmt.Printf("Unhandled server error: %s\n%s\n", err, string(debug.Stack()))
		w.WriteHeader(500)
		w.Write([]byte(http.StatusText(500)))
	}
}

// Creates a new HttpMux and initializes it.
// In most cases you should use this instead of router.New()
func NewHTTP() *HttpMux {
	m := new(HttpMux)
	m.WithDefaults()
	m.DefaultHost().ErrorHandler = DefaultErrorHandler
	return m
}

// Private helper method to serve up an HTTP error using the host error handler if applicable.
// Otherwise DefaultErrorHandler is used.
func (m *HttpMux) serveHTTPError(err error, w http.ResponseWriter, r *http.Request) {
	host := muxcontext.Host[http.Handler, HttpErrorHandler](r.Context())

	if host != nil && host.ErrorHandler != nil {
		host.ErrorHandler(err, w, r)
		return
	}

	defaultHost := m.DefaultHost()
	if defaultHost.ErrorHandler != nil {
		defaultHost.ErrorHandler(err, w, r)
		return
	}

	DefaultErrorHandler(err, w, r)
}

// ServeHTTP implements the standard HTTP interface can be used with most libraries that support HTTP handlers.
//
// This method modifies the request context:
//
//	`Host` will be the router.Host struct
//	`Resource` will be the router.Resource struct
//	`PathParams` will be the parameters parsed from the URL path
//	`HostParams` will be the parameters parsed from the hostname
func (m *HttpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() {
		if err := recover(); err != nil {
			switch errCast := err.(type) {
			case error:
				m.serveHTTPError(errCast, w, r.WithContext(ctx))
			default:
				fmt.Printf("%v\n", err)
				m.serveHTTPError(NewHttpError(http.StatusInternalServerError), w, r.WithContext(ctx))
			}
		}
	}()

	host, hostParamValues := m.Host(r.Host)
	resource, pathParamValues := host.Resource([]byte(r.URL.Path))

	ctx = muxcontext.WithHost(ctx, host)

	if resource == nil {
		m.serveHTTPError(NewHttpError(http.StatusNotFound), w, r.WithContext(ctx))
		return
	}

	ctx = muxcontext.WithResource(ctx, resource)

	// Normalize the method name to upper since this is be taken straight from the request header
	r.Method = strings.ToUpper(r.Method)

	handler, ok := resource.Method(r.Method)

	if ok {
		paramMap := resource.ParamMap(r.Method, pathParamValues)
		if paramMap != nil {
			ctx = muxcontext.WithPathParams(ctx, paramMap)
		}

		paramMap = host.ParamMap(hostParamValues)
		if paramMap != nil {
			ctx = muxcontext.WithHostParams(ctx, paramMap)
		}

		any(handler).(http.Handler).ServeHTTP(w, r.WithContext(ctx))
	} else if len(resource.Methods()) > 0 {
		// 405 Method Not Allowed - Has other methods but this isn't one
		m.serveHTTPError(NewHttpError(http.StatusMethodNotAllowed), w, r.WithContext(ctx))
	} else {
		// 404 Not Found - Has no methods at all
		m.serveHTTPError(NewHttpError(http.StatusNotFound), w, r.WithContext(ctx))
	}
}

func (m *HttpMux) HandleFunc(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(method, path, http.HandlerFunc(handler))
}
