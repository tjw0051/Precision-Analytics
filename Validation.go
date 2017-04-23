package main 

import(
	"strings"

	"precision-analytics/data"
)


func ValidateEntry(entry data.Entry) (bool, error) {

	//TODO:
	// - Validate UUIDs with regex
	//http://stackoverflow.com/questions/136505/searching-for-uuids-in-text-with-regex
	// - Validate/preformat types in msg
	// - Validate version number (format x.x.x.x OR flexible (e.g. 1.0.2a))
	// - Lowercase all required parameters
	entry.Id = strings.ToLower(entry.Id)
	entry.Platform = strings.ToLower(entry.Platform)
	entry.Namespace = strings.ToLower(entry.Namespace)
	entry.UserId = strings.ToLower(entry.UserId)
	entry.SessionId = strings.ToLower(entry.SessionId)
	// - Format date correctly - prevent dates earlier than server clock
	// - Char limit for all params
	return true, nil
}

// Placeholder func
// TODO: Read from prefs, allow/disallow non-auth entries
func RequiresAuth(entry data.Entry) bool {
	return true
}