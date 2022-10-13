package restkit

import (
	"fmt"
	"net/http"
	"strings"
)

type CORSHandler interface {
	HandleCORS(r *http.Request, policy *CORSHeaderPolicy) (requestedMethod string, realResponseHeaders map[string][]string)
}

type MethodRouting struct {
	PostHandler             http.Handler
	PutHandler              http.Handler
	DeleteHandler           http.Handler
	GetHandler              http.Handler
	PatchHandler            http.Handler
	ConnectHandler          http.Handler
	HeadHandler             http.Handler
	TraceHandler            http.Handler
	OptionsHandler          http.Handler
	MethodsCORSHeaderPolicy *MethodsCORSHeaderPolicy
	corsHeaderPolicy        *CORSHeaderPolicy
}

func (mr *MethodRouting) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && mr.PostHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodPost, &w)
		mr.PostHandler.ServeHTTP(w, r)
	case r.Method == http.MethodPut && mr.PutHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodPut, &w)
		mr.PutHandler.ServeHTTP(w, r)
	case r.Method == http.MethodDelete && mr.DeleteHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodDelete, &w)
		mr.DeleteHandler.ServeHTTP(w, r)
	case r.Method == http.MethodGet && mr.GetHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodGet, &w)
		mr.GetHandler.ServeHTTP(w, r)
	case r.Method == http.MethodPatch && mr.PatchHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodPatch, &w)
		mr.PatchHandler.ServeHTTP(w, r)
	case r.Method == http.MethodConnect && mr.ConnectHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodConnect, &w)
		mr.ConnectHandler.ServeHTTP(w, r)
	case r.Method == http.MethodHead && mr.HeadHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodHead, &w)
		mr.HeadHandler.ServeHTTP(w, r)
	case r.Method == http.MethodTrace && mr.TraceHandler != nil:
		mr.setCorsResponseHeaderIfExists(http.MethodTrace, &w)
		mr.TraceHandler.ServeHTTP(w, r)
	case r.Method == http.MethodOptions && mr.MethodsCORSHeaderPolicy != nil:
		corsResponseHeader, statusCode := (*mr.MethodsCORSHeaderPolicy).Validate(r)
		if corsResponseHeader == nil && statusCode < 200 || statusCode > 300 {
			w.WriteHeader(statusCode)
		} else {
			mr.corsHeaderPolicy = corsResponseHeader
			w.Header().Set(AccessControlAllowOrigin, corsResponseHeader.AccessControlAllowOrigin)
			w.Header().Set(AccessControlAllowCredentials, fmt.Sprintf("%t", corsResponseHeader.AccessControlAllowCredentials))
			w.Header().Set(AccessControlAllowMethods, strings.Join(*corsResponseHeader.AccessControlAllowMethods, ","))
			w.Header().Set(AccessControlAllowHeaders, strings.Join(*corsResponseHeader.AccessControlAllowHeaders, ","))
			w.WriteHeader(statusCode)
		}
	case r.Method == http.MethodOptions && mr.OptionsHandler != nil:
		mr.OptionsHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (mr *MethodRouting) setCorsResponseHeaderIfExists(method string, w *http.ResponseWriter) {
	if mr.corsHeaderPolicy != nil && mr.corsHeaderPolicy.method == method {
		defer func() { mr.corsHeaderPolicy = nil }()
		(*w).Header().Set(AccessControlAllowOrigin, mr.corsHeaderPolicy.AccessControlAllowOrigin)
		(*w).Header().Set(AccessControlAllowCredentials, fmt.Sprintf("%t", mr.corsHeaderPolicy.AccessControlAllowCredentials))
		(*w).Header().Set(AccessControlAllowMethods, strings.Join(*mr.corsHeaderPolicy.AccessControlAllowMethods, ","))
		(*w).Header().Set(AccessControlAllowHeaders, strings.Join(*mr.corsHeaderPolicy.AccessControlAllowHeaders, ","))
	}
}
