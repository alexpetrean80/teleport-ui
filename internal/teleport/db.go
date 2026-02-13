package teleport

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type TeleportDBLabels struct {
	CloudProvider string `json:"cloud-provider,omitempty"`
	DBName        string `json:"db-name,omitempty"`
	Engine        string `json:"engine,omitempty"`
	Identifier    string `json:"identifier,omitempty"`
	Owner         string `json:"owner,omitempty"`
}

type TeleportDBMeta struct {
	Name   string           `json:"name,omitempty"`
	Labels TeleportDBLabels `json:"labels"`
}

type TeleportUser string

func (u TeleportUser) String() string {
	return string(u)
}

type TeleportUsersList struct {
	Allowed []TeleportUser `json:"allowed,omitempty"`
}

type TeleportDB struct {
	Metadata TeleportDBMeta    `json:"metadata"`
	Users    TeleportUsersList `json:"users"`
}

func (db TeleportDB) String() string {
	engine := db.Metadata.Labels.Engine
	if engine == "" {
		engine = "unknown"
	}
	owner := db.Metadata.Labels.Owner
	if owner == "" {
		owner = "unknown"
	}
	return fmt.Sprintf("%s (%s) - owner: %s", db.Metadata.Name, engine, owner)
}

func GetTeleportDatabases(ctx context.Context) ([]TeleportDB, error) {
	cmd := exec.CommandContext(ctx, "tsh", "db", "ls", "--format", "json")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("creating stdout pipe: %s", err.Error())
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting command: %s", err.Error())
	}

	var dbs []TeleportDB

	if err := json.NewDecoder(stdout).Decode(&dbs); err != nil {
		return nil, fmt.Errorf("decoding databases: %s", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("waiting on cmd to finish: %s", err.Error())
	}

	return dbs, nil
}

func ConnectToTeleportDB(ctx context.Context, db *TeleportDB, user TeleportUser) error {
	cmd := exec.CommandContext(
		ctx,
		"tsh",
		"db",
		"connect",
		db.Metadata.Name,
		"--db-user",
		user.String(),
		"--db-name",
		db.Metadata.Labels.DBName,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
