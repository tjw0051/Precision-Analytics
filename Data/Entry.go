package data 

import "time"

type EntryMsg struct {
	// Key for this log entry
	Key 		string 		`json:"key"`
	// Supported Datatypes:
	// - bool
	// - string
	// - int
	// - float
	Type 		string 		`json:"type"`
	// Value for entry
	Value		string		`json:"value"`
}

type Entry struct {
	// JWT token
	Token 		string 		`json:"token"`
	// Unique Log Entry ID
	Id			string		`json:"id"`
	// Web, Android, iOS, Windows, MacOS
	// or comma-seperated list for multiple
	Platform	string		`json:"platform"`
	// App Namespace e.g. com.hypnosstudios.carkit
	Namespace	string		`json:"namespace"`
	// App Version
	Version 	string 		`json:"version"`
	// Unique User ID
	UserId		string		`json:"userId"`
	// Unique Session ID
	SessionId	string		`json:"sessionId"`
	// Date Entry was received
	Date		time.Time 	`json:"date"`
	// Categorises this set of entry data
	MsgType		string		`json:"msgType"`
	// JSON string of data (key-type-value format)
	Msg 		[]EntryMsg	`json:"msg"`
}

type Entries struct {
	Entries 	[]Entry 	`json:"entries"`
}

//type Entries []Entry