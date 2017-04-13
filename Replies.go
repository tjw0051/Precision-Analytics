package main

type ErrorReply struct {
	Errors 		[]ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Domain 		string 		`json:"domain"`
	Reason 		string 		`json:"reason"`
	Message 	string 		`json:"message"`
	Code 		string 		`json:"code"`
}

func GetError(domain string, reason string, message string, code string) ErrorReply {
	item := ErrorItem{Domain: domain, Reason: reason, Message: message, Code: code}
	reply := ErrorReply{Errors: []ErrorItem{item} }
	return reply
}