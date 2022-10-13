package restkit

import (
	"net/http"
	"strings"
)

const (
	Origin                        = "Origin"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
)

const (
	LowerCase uint8 = 0
	UpperCase uint8 = 1
)

type CaseSensitiveStrings *[]string

func NewCaseSensitiveStrings(lowerOrUpper uint8, strs ...string) CaseSensitiveStrings {
	var result []string

	for _, str := range strs {
		if lowerOrUpper == LowerCase {
			result = append(result, strings.ToLower(str))
		} else if lowerOrUpper == UpperCase {
			result = append(result, strings.ToUpper(str))
		} else {
			result = append(result, str)
		}
	}
	return &result
}

type CORSHeaderPolicy struct {
	AccessControlAllowOrigin      string
	AccessControlAllowCredentials bool
	AccessControlAllowMethods     CaseSensitiveStrings
	AccessControlAllowHeaders     CaseSensitiveStrings
	method                        string
}

type MethodsCORSHeaderPolicy map[string]CORSHeaderPolicy

func (mchp *MethodsCORSHeaderPolicy) Validate(r *http.Request) (*CORSHeaderPolicy, int) {
	// first check if the cors header policy itself is exists or not based on the http method
	// then check if the http method is allowed or not, if the requested method is allowed
	// then set the real method to the header policy so that when the real request is sent
	// the response header can be set correctly based on the method
	acReqMethod := r.Header.Get(AccessControlRequestMethod)
	if acReqMethod == "" {
		return nil, http.StatusBadRequest
	}
	corsHeaderPolicy, policyOk := (*mchp)[acReqMethod]
	if !policyOk {
		return nil, http.StatusMethodNotAllowed
	}
	corsHeaderPolicy.method = acReqMethod
	(*mchp)[acReqMethod] = corsHeaderPolicy

	accessControlMethodOk := false
	for _, allowMethod := range *corsHeaderPolicy.AccessControlAllowMethods {
		if acReqMethod == allowMethod {
			accessControlMethodOk = true
			break
		}
	}
	if !accessControlMethodOk {
		return nil, http.StatusMethodNotAllowed
	}

	reqOrigin := r.Header.Get(Origin)
	if reqOrigin == "" {
		return nil, http.StatusBadRequest
	}
	if reqOrigin != corsHeaderPolicy.AccessControlAllowOrigin {
		return nil, http.StatusBadRequest
	}

	acReqHeaderRaw := strings.ToLower(r.Header.Get(AccessControlRequestHeaders))
	if acReqHeaderRaw == "" {
		return nil, http.StatusBadRequest
	}
	acReqHeaders := strings.Split(strings.ReplaceAll(acReqHeaderRaw, " ", ""), ",")

	accessControlHeadersOk := true
	var innerAccessControlHeadersOk bool
	for _, allowHeader := range *corsHeaderPolicy.AccessControlAllowHeaders {
		innerAccessControlHeadersOk = false
		for _, acReqHeader := range acReqHeaders {
			if acReqHeader == allowHeader {
				innerAccessControlHeadersOk = true
				break
			}
		}
		if !innerAccessControlHeadersOk {
			accessControlHeadersOk = false
			break
		}
	}
	if !accessControlHeadersOk {
		return nil, http.StatusBadRequest
	}

	return &corsHeaderPolicy, http.StatusNoContent
}
