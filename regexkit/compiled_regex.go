package regexkit

import (
	"regexp"

	"github.com/ilhammhdd/go-toolkit/errorkit"
)

const (
	RegexEmail uint = iota
	RegexAlphanumeric
	RegexNotEmpty
	RegexURL
	RegexJWT
	RegexNumber
	RegexLatitude
	RegexLongitude
	RegexUUIDV4
	RegexCommonUnitOfLength
	RegexIPv4
	RegexIPv4TCPPortRange
	RegexDateTimeRFC3339
	ParamNotExists
	LastRegexIota
)

var Regex map[uint]string = map[uint]string{
	RegexEmail:              `^[A-Za-z0-9](([_\.\-]?[a-zA-Z0-9]+)*)@([A-Za-z0-9]+)(([\.\-]?[a-zA-Z0-9]+)*)\.([A-Za-z]{2,})$`,
	RegexAlphanumeric:       `^[a-zA-Z0-9]+$`,
	RegexNotEmpty:           `.*\S.*`,
	RegexURL:                `^((((h)(t)|(f))(t)(p)((s)?))\://)?(www.|[a-zA-Z0-9].)[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,6}(\:[0-9]{1,5})*(/($|[a-zA-Z0-9\.\,\;\?\'\\\+&amp;%\$#\=~_\-]+))*$`,
	RegexJWT:                `^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`,
	RegexNumber:             `^[0-9]+`,
	RegexLatitude:           `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`,
	RegexLongitude:          `^[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`,
	RegexUUIDV4:             `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`,
	RegexCommonUnitOfLength: `^(mm|MM|cm|CM|dm|DM|m|M|dam|DAM|hm|HM|km|KM)$`,
	RegexIPv4:               `^(([0-9]{1,2}|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]{1,2}|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`,
	RegexIPv4TCPPortRange:   `^(([0-9]{1,2}|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]{1,2}|1[0-9]{2}|2[0-4][0-9]|25[0-5])\:([0-9]{1,4}|[1-5][0-9]{1,4}|6553[0-5]|655[0-2][0-9]|65[0-4][0-9][0-9]|6[0-4][0-9][0-9][0-9])$`,
	RegexDateTimeRFC3339:    `^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`,
}

var RegexpCompiled map[uint]*regexp.Regexp

func CompileAllRegex(otherRegex map[uint]string) error {
	RegexpCompiled = make(map[uint]*regexp.Regexp)
	var err error

	regexEmailComp, err := regexp.Compile(Regex[RegexEmail])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexEmail] = regexEmailComp

	regexAlphanumericComp, err := regexp.Compile(Regex[RegexAlphanumeric])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexAlphanumeric] = regexAlphanumericComp

	regexNotEmptyComp, err := regexp.Compile(Regex[RegexNotEmpty])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexNotEmpty] = regexNotEmptyComp

	regexURLComp, err := regexp.Compile(Regex[RegexURL])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexURL] = regexURLComp

	regexJWTComp, err := regexp.Compile(Regex[RegexJWT])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexJWT] = regexJWTComp

	regexNumberComp, err := regexp.Compile(Regex[RegexNumber])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexNumber] = regexNumberComp

	regexLatitudeComp, err := regexp.Compile(Regex[RegexLatitude])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexLatitude] = regexLatitudeComp

	regexLongitudeComp, err := regexp.Compile(Regex[RegexLongitude])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexLongitude] = regexLongitudeComp

	regexUUIDV4Comp, err := regexp.Compile(Regex[RegexUUIDV4])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexUUIDV4] = regexUUIDV4Comp

	regexCommonUnitOfLengthComp, err := regexp.Compile(Regex[RegexCommonUnitOfLength])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexCommonUnitOfLength] = regexCommonUnitOfLengthComp

	regexIPv4Comp, err := regexp.Compile(Regex[RegexIPv4])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexIPv4] = regexIPv4Comp

	regexIPv4TCPPortRangeComp, err := regexp.Compile(Regex[RegexIPv4TCPPortRange])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexIPv4TCPPortRange] = regexIPv4TCPPortRangeComp

	regexDateTimeRFC3339Comp, err := regexp.Compile(Regex[RegexDateTimeRFC3339])
	if errorkit.ErrorHandled(err, stackSize) {
		return err
	}
	RegexpCompiled[RegexDateTimeRFC3339] = regexDateTimeRFC3339Comp

	for key, val := range otherRegex {
		if key < LastRegexIota {
			continue
		}
		regexCompiled, err := regexp.Compile(val)
		errorkit.ErrorHandled(err, stackSize)
		RegexpCompiled[key] = regexCompiled
	}

	return nil
}
