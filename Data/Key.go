package data 

import (
	"time"
)

type Keys struct {
	Keys 			[]Key 	`json:"keys"`
}

type Key struct {
	Key 			string 		`json:"key"`
	Expires 	 	bool 		`json:"expires"`
	ExpDate		 	time.Time 	`json:"expDate"`
	Active	 		bool 		`json:"active"`
	Group 		  	string 		`json:"group"`
}