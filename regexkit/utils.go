package regexkit

const stackSize = 1024 * 8

type RegexNoMatchMsgGenerator interface {
	GenerateRegexNoMatchMsg(regexConst int) (noMatchMsg string)
	GenerateParamNotExistsMsg(param string) (noMatchMsg string)
}
