package data 

type Groups struct {
	Groups 			[]Group 	`json:"groups"`
}

type Group struct {
	Name 			string 		`json:"name"`
	ManageGroups 	bool 		`json:"manageGroups"`
	ManageServer 	bool 		`json:"manageServer"`
	CreateKey 		bool 		`json:"createKey"`
	AddLog 			bool 		`json:"addLog"`
	QueryLog 		bool 		`json:"queryLog"`
	ModLog 			bool 		`json:"modLog"`
}