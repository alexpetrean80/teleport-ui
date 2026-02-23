package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexptr80/teleport-ui/internal/teleport"
	"github.com/alexptr80/teleport-ui/internal/ui"
)

func main() {
	ctx := context.Background()
	dbs, err := teleport.GetTeleportDatabases(ctx)
	if err != nil {
		log.Fatal(err)
	}

	selectedDB, err := ui.RunFuzzyFinder(dbs)
	if err != nil {
		log.Fatal(err)
	}

	if selectedDB == nil {
		fmt.Println("No database selected")
		return
	}

	fmt.Printf("Selected: %s\n", selectedDB.String())

	selectedDBUser, err := ui.RunFuzzyFinder(selectedDB.Users.Allowed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("selected %s", selectedDBUser)
	if err = teleport.ConnectToTeleportDB(ctx, selectedDB, *selectedDBUser); err != nil {
		log.Fatal(err)
	}
}
