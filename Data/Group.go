package data 

type Groups struct {
	Groups 			[]Group 	`json:"groups"`
}

type Group struct {
	Name 			string 		`json:"name"`
	Perms		 	string 		`json:"perms"`
}