package teleport

import (
	"encoding/json"
	"testing"
)

func TestTeleportDBString(t *testing.T) {
	tests := []struct {
		name string
		db   TeleportDB
		want string
	}{
		{
			name: "all labels present",
			db: TeleportDB{
				Metadata: TeleportDBMeta{
					Name: "prod-db",
					Labels: TeleportDBLabels{
						Engine: "postgres",
						Owner:  "platform",
					},
				},
			},
			want: "prod-db (postgres) - owner: platform",
		},
		{
			name: "missing engine",
			db: TeleportDB{
				Metadata: TeleportDBMeta{
					Name: "staging-db",
					Labels: TeleportDBLabels{
						Owner: "backend",
					},
				},
			},
			want: "staging-db (unknown) - owner: backend",
		},
		{
			name: "missing owner",
			db: TeleportDB{
				Metadata: TeleportDBMeta{
					Name: "dev-db",
					Labels: TeleportDBLabels{
						Engine: "mysql",
					},
				},
			},
			want: "dev-db (mysql) - owner: unknown",
		},
		{
			name: "both missing",
			db: TeleportDB{
				Metadata: TeleportDBMeta{
					Name: "bare-db",
				},
			},
			want: "bare-db (unknown) - owner: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.db.String()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTeleportKubeClusterString(t *testing.T) {
	tests := []struct {
		name    string
		cluster TeleportKubeCluster
		want    string
	}{
		{
			name: "with cloudProvider",
			cluster: TeleportKubeCluster{
				KubeClusterName: "prod-k8s",
				Labels:          TeleportKubeLabels{"cloudProvider": "aws"},
			},
			want: "prod-k8s (aws)",
		},
		{
			name: "with clusterName only",
			cluster: TeleportKubeCluster{
				KubeClusterName: "staging-k8s",
				Labels:          TeleportKubeLabels{"clusterName": "staging-main"},
			},
			want: "staging-k8s (staging-main)",
		},
		{
			name: "cloudProvider takes precedence",
			cluster: TeleportKubeCluster{
				KubeClusterName: "multi-k8s",
				Labels: TeleportKubeLabels{
					"cloudProvider": "gcp",
					"clusterName":   "my-cluster",
				},
			},
			want: "multi-k8s (gcp)",
		},
		{
			name: "no labels",
			cluster: TeleportKubeCluster{
				KubeClusterName: "bare-k8s",
			},
			want: "bare-k8s",
		},
		{
			name: "irrelevant labels",
			cluster: TeleportKubeCluster{
				KubeClusterName: "other-k8s",
				Labels:          TeleportKubeLabels{"env": "prod"},
			},
			want: "other-k8s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cluster.String()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTeleportUserString(t *testing.T) {
	u := TeleportUser("admin")
	if got := u.String(); got != "admin" {
		t.Errorf("got %q, want %q", got, "admin")
	}
}

func TestDBJSONUnmarshal(t *testing.T) {
	raw := `[
		{
			"metadata": {
				"name": "my-postgres",
				"labels": {
					"engine": "postgres",
					"owner": "team-a",
					"db-name": "appdb",
					"cloud-provider": "aws",
					"identifier": "rds-123"
				}
			},
			"users": {
				"allowed": ["readonly", "admin"]
			}
		}
	]`

	var dbs []TeleportDB
	if err := json.Unmarshal([]byte(raw), &dbs); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(dbs) != 1 {
		t.Fatalf("got %d dbs, want 1", len(dbs))
	}

	db := dbs[0]
	if db.Metadata.Name != "my-postgres" {
		t.Errorf("Name = %q, want %q", db.Metadata.Name, "my-postgres")
	}
	if db.Metadata.Labels.Engine != "postgres" {
		t.Errorf("Engine = %q, want %q", db.Metadata.Labels.Engine, "postgres")
	}
	if db.Metadata.Labels.CloudProvider != "aws" {
		t.Errorf("CloudProvider = %q, want %q", db.Metadata.Labels.CloudProvider, "aws")
	}
	if db.Metadata.Labels.DBName != "appdb" {
		t.Errorf("DBName = %q, want %q", db.Metadata.Labels.DBName, "appdb")
	}
	if len(db.Users.Allowed) != 2 {
		t.Fatalf("got %d users, want 2", len(db.Users.Allowed))
	}
	if db.Users.Allowed[0] != "readonly" {
		t.Errorf("Users[0] = %q, want %q", db.Users.Allowed[0], "readonly")
	}
}

func TestKubeJSONUnmarshal(t *testing.T) {
	raw := `[
		{
			"kube_cluster_name": "prod-cluster",
			"labels": {
				"cloudProvider": "gcp",
				"env": "production"
			},
			"selected": true
		}
	]`

	var clusters []TeleportKubeCluster
	if err := json.Unmarshal([]byte(raw), &clusters); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(clusters) != 1 {
		t.Fatalf("got %d clusters, want 1", len(clusters))
	}

	c := clusters[0]
	if c.KubeClusterName != "prod-cluster" {
		t.Errorf("KubeClusterName = %q, want %q", c.KubeClusterName, "prod-cluster")
	}
	if !c.Selected {
		t.Error("Selected = false, want true")
	}
	if c.Labels["cloudProvider"] != "gcp" {
		t.Errorf("Labels[cloudProvider] = %q, want %q", c.Labels["cloudProvider"], "gcp")
	}
}
