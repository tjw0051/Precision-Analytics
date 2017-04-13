package db 

import(
	//"fmt"
	"log"

	"Precision-Analytics/data"

	"github.com/nu7hatch/gouuid"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	entries := []data.EntryMsg { data.EntryMsg{Key: "1", Type: "t1", Value: "v1"},
	 							data.EntryMsg{Key: "2", Type: "t2", Value: "v2"} }

	Add(data.Entry{Id: "0000000002", Msg: entries})
}

func Create() {
	log.Printf("Creating Database: PA.db")
	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)
	defer db.Close()

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
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func Add(entry data.Entry) {

	id, err := uuid.NewV4()
	checkErr(err)
	entry.Id = id.String()

	db, err := sql.Open("sqlite3", "./PA.db")
	checkErr(err)

	stmt, err := db.Prepare("INSERT INTO log(id, platform, namespace, version, userId, sessionId, date, msgType, key, type, value) values(?,?,?,?,?,?,?,?,?,?,?)")
    checkErr(err)

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
    db.Close()
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


