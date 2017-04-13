package data

type ErrorReply struct {
	Errors 		[]ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Domain 		string 		`json:"domain"`
	Reason 		string 		`json:"reason"`
	Message 	string 		`json:"message"`
	Code 		int 		`json:"code"`
}

var Errors map[string]ErrorItem = map[string]ErrorItem { 
	"jsonInvalid": {"msgFormat","jsonInvalid","JSON could not be parsed.",422},
	"keyInvalid": {"auth","keyInvalid","API Key was invalid.",401},
	"tokenInvalid": {"auth","tokenInvalid","API Token was invalid.",401},
	"tokenErr": {"auth","tokenErr","Could not create token.",401} }

func GetError(domain string, reason string, message string, code int) ErrorReply {
	item := ErrorItem{Domain: domain, Reason: reason, Message: message, Code: code}
	reply := ErrorReply{Errors: []ErrorItem{item} }
	return reply
}