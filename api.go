package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"

	"signinapc.app/api/pkg/endpoints"
	"signinapc.app/api/pkg/events"
	"signinapc.app/api/pkg/services"
)

func main() {
	fmt.Println("Hello!")

	ids, err := events.MakeRegexIDs()

	if err != nil {
		fmt.Println("Failed to create ID regex:", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", "./test.db")

	if err != nil {
		fmt.Println("Failed to open database:", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	_, err = db.Exec("SELECT 1 FROM events")

	if err != nil {
		fmt.Println("Failed to query events table:", err.Error())
		os.Exit(1)
	}

	events := services.MakeSqliteEventService(db, ids)

	memberSignIn := endpoints.MemberSignIn(events, ids)
	memberAgreed := endpoints.MemberAgreed(events, ids)

	mux := http.NewServeMux()
	mux.HandleFunc(memberSignIn.Path(), memberSignIn.Handle)
	mux.HandleFunc(memberAgreed.Path(), memberAgreed.Handle)

	err = http.ListenAndServe(":3339", mux)

	if err != nil {
		fmt.Println("Failed to start server.", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
