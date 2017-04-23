package data 

import(
	"fmt"
	"log"
	"time"

	//"Precision-Analytics/data"

	"github.com/nu7hatch/gouuid"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	entries := []EntryMsg { EntryMsg{Key: "1", Type: "t1", Value: "v1"},
	 							EntryMsg{Key: "2", Type: "t2", Value: "v2"} }
	LoadConfig("pa.config")
	CreateDB()

	SetLog(Entry{Id: "0000000002", Msg: entries})
}

/*******************************************************
	Database: PA.db
	Tables:
		log
		--------
		Stores analytics messages received from clients
		- id 			- Not DB unique - see docs
		- platform 		- e.g web, android, ios
		- namespace 	- e.g. com.companyname.appname
		- version 		- app version, e.g. 1.0.2
		- userId 		- uniquely identify user
		- sessionId 	- uniquely identify session
		- date 			- date message was received
		- msgType 		- user-defined msg type
		- key 			- parameter key
		- type 			- parameter type
		- value 		- parameter value

		apikeys
		--------
		Provide access to the API for requesting tokens
		- key 			- The API Key
		- expires 		- Does key expire or not
		- expDate 		- expiration date IF expires == true
		- active 		- active/deactivate a key without removing
		- accessGroup 	- access group (what the key can do)

		accessgroups
		-------
		Defines groups with different levels of access to the
		API. API keys and users are members of groups.
		- name 			- Name of group
		- manageGroups	- Can set groups
		- manageServer 	- Can manage (start/stop) server
		- createKey 	- Can create keys for own group
		- addLog 		- Can log messages
		- queryLog 		- Can get results from log
		- modLog 		- Can modify/remove from log

*******************************************************/

func CreateDB() {
	log.Printf("Creating Database: PA.db")
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()
	// Create Log Table
	sqlStmt := `
	create table log (
		id text not null,
		platform text,
		namespace text,
		version text,
		userId text,
		sessionId text,
		date timestamp,
		msgType text,
		key text,
		type text,
		value text);
	delete from log;
	`
	log.Printf("Creating Table: log")
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	// Create apikeys table
	sqlStmt = `
	create table apikeys (
		key text not null,
		expires boolean not null,
		expDate timestamp,
		active boolean not null,
		accessGroup text not null);
	delete from apikeys;
	`
	log.Printf("Creating Table: apikeys")
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	// Create groups table
	sqlStmt = `
	create table accessgroups (
		name text not null unique,
		perms text);
	delete from accessgroups;
	`
	log.Printf("Creating Table: accessgroups")
	_, err = db.Exec(sqlStmt)
	checkErr(err)

	// Create root group
	// Permissions is comma-delimited list of route names. * = all
	rootGroup := Group{
		Name: "root",
		Perms: "*",
	}
	// Wildcard group for all users
	allGroup := Group{
		Name: "*",
		Perms: "ReqAuth",
	}
	// Basic group for logging messages from application
	standardClientGroup := Group{
		Name: "client",
		Perms: "ReqAuth,Log",
	}

	rootGroups := Groups{Groups: []Group{rootGroup, allGroup, standardClientGroup}}
	SetGroups(rootGroups)

	// Create root API Key from config
	rootKey := Key{
		Key: ROOTKEY,
		Expires: false,
		ExpDate: time.Now(),
		Active: true,
		Group: "root",
	}
	rootKeys := Keys{Keys: []Key{rootKey}}
	SetKeys(rootKeys)

	//TODO create root group and root key from config

}

func SetLog(entry Entry) {

	id, err := uuid.NewV4()
	checkErr(err)
	entry.Id = id.String()

	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("INSERT INTO log(id, platform, namespace, version, userId, sessionId, date, msgType, key, type, value) values(?,?,?,?,?,?,?,?,?,?,?)")
    checkErr(err)
    defer stmt.Close()

    // TODO: Split up entry into seperate entries

    for i := 0; i < len(entry.Msg); i++ {
    	_, err = stmt.Exec(entry.Id, 
    	entry.Platform, 
    	entry.Namespace, 
    	entry.Version, 
    	entry.UserId,
    	entry.SessionId,
    	entry.Date,
    	entry.MsgType,
    	entry.Msg[i].Key,
    	entry.Msg[i].Type,
    	entry.Msg[i].Value)
    	checkErr(err)
    }
    tx.Commit()
}

func GetLog() Entries {
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("select id, platform, namespace, version, userId, sessionId, date, msgType, key, type, value from log order by id asc")
	checkErr(err)
	defer rows.Close()

	entries := Entries{}

	for rows.Next() {
		e := Entry{Msg: []EntryMsg { EntryMsg{}} }
		err = rows.Scan(&e.Id, &e.Platform, &e.Namespace, &e.Version, &e.UserId, &e.SessionId, &e.Date, &e.MsgType, &e.Msg[0].Key, &e.Msg[0].Type, &e.Msg[0].Value)
		if(len(entries.Entries) > 0) {
			if(entries.Entries[0].Id == e.Id) {
				entries.Entries[0].Msg = append(entries.Entries[0].Msg, e.Msg[0])
			} else {
				entries.Entries = append(entries.Entries, e)
			}
		}
	}
	err = rows.Err()
	checkErr(err)

	return entries
}

func GetGroups() Groups {
	// TODO: Cache groups
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("select name, perms from accessgroups")
	checkErr(err)
	defer rows.Close()

	groups := Groups{}

	for rows.Next() {
		newGroup := Group{}
		err = rows.Scan(&newGroup.Name,
			&newGroup.Perms)
		checkErr(err)
		groups.Groups = append(groups.Groups, newGroup)
	}
	err = rows.Err()
	checkErr(err)
	return groups
}

func GetGroup(name string) (Group, error) {
	groups := GetGroups()
	for i := 0; i < len(groups.Groups); i++ {
		if(groups.Groups[i].Name == name) {
			return groups.Groups[i], nil
		}
	}
	return Group{}, fmt.Errorf("Group not found.")
}

func SetGroups(groups Groups) {
	//TODO: Check if group exists -  if it does, overwrite
	//TODO: DO NOT allow changes to root group
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("INSERT INTO accessgroups(name, perms) values(?,?)")
	checkErr(err)
	defer stmt.Close()

	for i := 0; i < len(groups.Groups); i++ {
		_, err = stmt.Exec(groups.Groups[i].Name,
			groups.Groups[i].Perms)
		checkErr(err)
	}
	tx.Commit()
}

func RemoveGroups(names []string) {
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("delete from accessgroups where name = ?")
	checkErr(err)
	defer stmt.Close()

	for i := 0; i < len(names); i++ {
		if(names[i] == "root") {
			// DO NOT remove root group
		}
		_, err := stmt.Exec(names[i])
		checkErr(err)
	}

	//TODO: Revoke keys tied to group
}

func GetKeys() Keys {
	// TODO: Cache keys
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("select key, expires, expDate, active, accessGroup from apikeys")
	checkErr(err)
	defer rows.Close()

	keys := Keys{}

	for rows.Next() {
		newKey := Key{}
		err = rows.Scan(&newKey.Key,
						 &newKey.Expires,
						 &newKey.ExpDate,
						 &newKey.Active,
						 &newKey.Group)
		checkErr(err)
		keys.Keys = append(keys.Keys, newKey)
	}
	err = rows.Err()
	checkErr(err)

	return keys
}

func SetKeys(keys Keys) Keys {
	// TODO: Check if key exists. If it does, overwrite
	// 		 Check if group exists
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("INSERT INTO apikeys(key, expires, expDate, active, accessGroup) values(?,?,?,?,?)")
	checkErr(err)
	defer stmt.Close()

	for i := 0; i < len(keys.Keys); i++ {
		// Generate Key
		if(keys.Keys[i].Key == "") {
			id, err := uuid.NewV4()
			checkErr(err)
			keys.Keys[i].Key = id.String()
		}
		_, err = stmt.Exec(keys.Keys[i].Key,
			keys.Keys[i].Expires,
			keys.Keys[i].ExpDate,
			keys.Keys[i].Active,
			keys.Keys[i].Group)
		checkErr(err)
	}
	tx.Commit()

	return keys
}

func RemoveKeys(keys []string) {
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare("delete from apikeys where key = ?")
	checkErr(err)
	defer stmt.Close()

	for i := 0; i < len(keys); i++ {
		_, err := stmt.Exec(keys[i])
		checkErr(err)
	}
}

// Get(id) - Find all with ID, construct Entry and return

// Get(query) - Run SQL query and return results

// Get(id, platform, namespace, version, userId, sessionId, date, msgType)
// - Use filters to return all matching values
// - TODO: API pagination with tokens (like Google's API)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}


