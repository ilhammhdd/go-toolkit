package restkit

import (
	"net/http"
	"net/url"

	"github.com/ilhammhdd/go-toolkit/regexkit"
)

type HeaderParamValidation struct {
	RegexRules map[string]int
	Header     http.Header
}

func (hpv *HeaderParamValidation) Validate(errMsgGen regexkit.RegexNoMatchMsgGenerator) (*map[string][]string, bool) {
	valid := true
	var allErrs map[string][]string = make(map[string][]string)

	for param, regexConst := range hpv.RegexRules {
		if values, ok := hpv.Header[param]; !ok {
			valid = false
			allErrs[param] = []string{errMsgGen.GenerateParamNotExistsMsg(param)}
		} else {
			errs := make([]string, len(values))
			for idx, value := range values {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(value))
				if !regexOk {
					valid = false
					errs[idx] = errMsgGen.GenerateRegexNoMatchMsg(regexConst)
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

type URLParamValidation struct {
	RegexRules map[string]int
	Values     url.Values
}

func (upv *URLParamValidation) Validate(errMsgGen URLParamValidationRegexMsgGenerator) (*map[string][]string, bool) {
	var allErrs map[string][]string = make(map[string][]string)

	if len(upv.Values) == 0 {
		allErrs["all"] = []string{errMsgGen.GenerateNoURLParamMsg()}
		return &allErrs, false
	}

	valid := true

	for param, regexConst := range upv.RegexRules {
		vals, ok := upv.Values[param]
		if !ok && len(upv.Values) > 0 {
			valid = false
			allErrs[param] = []string{errMsgGen.GenerateParamNotExistsMsg(param)}
		} else {
			errs := make([]string, len(vals))
			for idx, val := range vals {
				regexOk := regexkit.RegexpCompiled[regexConst].Match([]byte(val))
				if !regexOk {
					valid = false
					errs[idx] = errMsgGen.GenerateRegexNoMatchMsg(regexConst)
				}
			}
			allErrs[param] = errs
		}
	}

	return &allErrs, valid
}
