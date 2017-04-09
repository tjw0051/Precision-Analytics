package main 

import(
	"strings"
)


func ValidateEntry(entry Entry) (bool, error) {

	if(RequiresAuth(entry)) {
		// Check token is valid
		if valid, err := ValidateToken(entry.Token); !valid {
			return valid, err
		}
	}
	//TODO:
	// - Validate UUIDs with regex
	// - Validate/preformat types in msg
	// - Lowercase all required parameters
	entry.Id = strings.ToLower(entry.Id)
	entry.Platform = strings.ToLower(entry.Platform)
	entry.Namespace = strings.ToLower(entry.Namespace)
	entry.userId = strings.ToLower(entry.userId)
	entry.sessionId = strings.ToLower(entry.sessionId)
	// - Format date correctly - prevent dates earlier than server clock
	// - Char limit for all params
	return true, nil
}

// Placeholder func
// TODO: Read from prefs, allow/disallow non-auth entries
func RequiresAuth(entry Entry) bool {
	return true
}