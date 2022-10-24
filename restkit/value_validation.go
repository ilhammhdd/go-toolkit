package restkit

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
)

type HeaderParamValidation struct {
	RegexRules map[string]uint
	Header     http.Header
}

func (hpv *HeaderParamValidation) Validate(detailedErrDescGen, regexErrDescGen errorkit.ErrDescGenerator) (*map[string][]string, bool) {
	valid := true
	var allErrs map[string][]string = make(map[string][]string)

	for param, regexConst := range hpv.RegexRules {
		if values, ok := hpv.Header[param]; !ok {
			valid = false
			allErrs[param] = []string{detailedErrDescGen.GenerateDesc(errorkit.FlowErrHttpHeaderParamNotExists, param)}
		} else {
			errs := make([]string, len(values))
			for idx, value := range values {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(value))
				if !regexOk {
					valid = false
					errs[idx] = regexErrDescGen.GenerateDesc(regexConst, param, fmt.Sprint(idx))
				}
			}
			allErrs[param] = errs
		}
	}

	return &allErrs, valid
}

type URLQueryValidation struct {
	RegexRules map[string]uint
	Values     url.Values
}

func (upv *URLQueryValidation) Validate(detailedErrDescGen, regexErrDescGen errorkit.ErrDescGenerator) (*map[string][]string, bool) {
	var allErrs map[string][]string = make(map[string][]string)

	if len(upv.Values) == 0 {
		allErrs["all"] = []string{detailedErrDescGen.GenerateDesc(errorkit.FlowErrURLQueryNotExists)}
		return &allErrs, false
	}

	valid := true

	for query, regexConst := range upv.RegexRules {
		vals, ok := upv.Values[query]
		if !ok && len(upv.Values) > 0 {
			valid = false
			allErrs[query] = []string{regexErrDescGen.GenerateDesc(regexConst, query)}
		} else {
			errs := make([]string, len(vals))
			for idx, val := range vals {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(val))
				if !regexOk {
					valid = false
					errs[idx] = regexErrDescGen.GenerateDesc(regexConst, query, fmt.Sprint(idx))
				}
			}
			allErrs[query] = errs
		}
	}

	return &allErrs, valid
}
