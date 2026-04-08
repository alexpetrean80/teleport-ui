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

type TeleportKubeLabels map[string]string

type TeleportKubeCluster struct {
	KubeClusterName string             `json:"kube_cluster_name"`
	Labels          TeleportKubeLabels `json:"labels"`
	Selected        bool               `json:"selected"`
}

func (c TeleportKubeCluster) String() string {
	labels := ""
	if len(c.Labels) > 0 {
		if provider, ok := c.Labels["cloudProvider"]; ok {
			labels = fmt.Sprintf(" (%s)", provider)
		}
		if name, ok := c.Labels["clusterName"]; ok && labels == "" {
			labels = fmt.Sprintf(" (%s)", name)
		}
	}
	return fmt.Sprintf("%s%s", c.KubeClusterName, labels)
}

const kubeCacheName = "kube"

// GetTeleportKubeClusters lists Kubernetes clusters the user can login to.
// When clearCache is false, cached results are returned if available.
// Filter args are passed through to tsh (e.g. key1=value1, --search=foo, --query='...').
func GetTeleportKubeClusters(ctx context.Context, filterArgs []string, clearCache bool) ([]TeleportKubeCluster, error) {
	if clearCache {
		if err := cache.Clear(kubeCacheName); err != nil {
			return nil, fmt.Errorf("clearing kube cache: %w", err)
		}
	}

	if !clearCache && len(filterArgs) == 0 {
		var cached []TeleportKubeCluster
		if err := cache.Load(kubeCacheName, &cached); err == nil {
			return cached, nil
		}
	}

	args := []string{"kube", "ls", "--format", "json"}
	args = append(args, filterArgs...)

	cmd := exec.CommandContext(ctx, "tsh", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running tsh kube ls: %s", err.Error())
	}

	var clusters []TeleportKubeCluster

	if stdout.Len() > 0 {
		if err := json.Unmarshal(stdout.Bytes(), &clusters); err != nil {
			return nil, fmt.Errorf("decoding kube clusters: %s", err.Error())
		}
	}

	if len(filterArgs) == 0 {
		_ = cache.Save(kubeCacheName, clusters)
	}

	return clusters, nil
}

// LoginToTeleportKube runs tsh kube login for the selected cluster.
func LoginToTeleportKube(ctx context.Context, cluster *TeleportKubeCluster) error {
	cmd := exec.CommandContext(ctx, "tsh", "kube", "login", cluster.KubeClusterName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
