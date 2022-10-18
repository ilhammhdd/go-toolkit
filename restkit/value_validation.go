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

func (hpv *HeaderParamValidation) Validate(errDescGen errorkit.ErrDescGenerator) (*map[string][]string, bool) {
	valid := true
	var allErrs map[string][]string = make(map[string][]string)

	for param, regexConst := range hpv.RegexRules {
		if values, ok := hpv.Header[param]; !ok {
			valid = false
			allErrs[param] = []string{errDescGen.GenerateDesc(regexkit.ParamNotExists, param)}
		} else {
			errs := make([]string, len(values))
			for idx, value := range values {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(value))
				if !regexOk {
					valid = false
					errs[idx] = errDescGen.GenerateDesc(regexConst, param, fmt.Sprint(idx))
				}
			}
			allErrs[param] = errs
		}
	}

	return &allErrs, valid
}

type URLParamValidationRegexMsgGenerator interface {
	regexkit.RegexNoMatchMsgGenerator
	GenerateNoURLParamMsg() string
}

const ErrNoURLParams uint = 0

type URLParamValidation struct {
	RegexRules map[string]uint
	Values     url.Values
}

func (upv *URLParamValidation) Validate(errDescGen errorkit.ErrDescGenerator) (*map[string][]string, bool) {
	var allErrs map[string][]string = make(map[string][]string)

	if len(upv.Values) == 0 {
		allErrs["all"] = []string{errDescGen.GenerateDesc(ErrNoURLParams)}
		return &allErrs, false
	}

	valid := true

	for param, regexConst := range upv.RegexRules {
		vals, ok := upv.Values[param]
		if !ok && len(upv.Values) > 0 {
			valid = false
			allErrs[param] = []string{errDescGen.GenerateDesc(regexConst, param)}
		} else {
			errs := make([]string, len(vals))
			for idx, val := range vals {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(val))
				if !regexOk {
					valid = false
					errs[idx] = errDescGen.GenerateDesc(regexConst, param, fmt.Sprint(idx))
				}
			}
			allErrs[param] = errs
		}
	}

	return &allErrs, valid
}
