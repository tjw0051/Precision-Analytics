// Mock Database for Tutorial
package main 

import "fmt"

var entries Entries

func init() {
	RepoCreateEntry(Entry{Id: "0000000000"})
	RepoCreateEntry(Entry{Id: "0000000001"})
}

func RepoFindEntry(id string) Entry {
	for _, t := range entries {
		if t.Id == id {
			return t
		}
	}
	return Entry{} // return empty Todo
}

func RepoCreateEntry(t Entry) Entry {
	entries = append(entries, t)
	return t
}

func RepoDestroyEntry(id string) error {
	for i, t := range entries {
		if t.Id == id {
			// Variadic Parameter: ...
			// Passes the rest of the set, from that point onward
			entries = append(entries[:i], entries[i+1:]...)
		}
	}
	return fmt.Errorf("Could not find Todo with id of %s to delete", id)
}


