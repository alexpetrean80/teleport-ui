package teleport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/alexpetrean80/teleport-ui/internal/cache"
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

const dbCacheName = "db"

// GetTeleportDatabases lists databases. When clearCache is false, cached results
// are returned if available. Filter args are passed through to tsh
// (e.g. key1=value1, --search=foo, --query='...').
func GetTeleportDatabases(ctx context.Context, filterArgs []string, clearCache bool) ([]TeleportDB, error) {
	if clearCache {
		if err := cache.Clear(dbCacheName); err != nil {
			return nil, fmt.Errorf("clearing db cache: %w", err)
		}
	}

	if !clearCache && len(filterArgs) == 0 {
		var cached []TeleportDB
		if err := cache.Load(dbCacheName, &cached); err == nil {
			return cached, nil
		}
	}

	args := []string{"db", "ls", "--format", "json"}
	args = append(args, filterArgs...)

	cmd := exec.CommandContext(ctx, "tsh", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running tsh db ls: %s", err.Error())
	}

	var dbs []TeleportDB

	if stdout.Len() > 0 {
		if err := json.Unmarshal(stdout.Bytes(), &dbs); err != nil {
			return nil, fmt.Errorf("decoding databases: %s", err.Error())
		}
	}

	if len(filterArgs) == 0 {
		_ = cache.Save(dbCacheName, dbs)
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
