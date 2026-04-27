package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alexpetrean80/teleport-ui/internal/cache"
	"github.com/alexpetrean80/teleport-ui/internal/teleport"
	"github.com/alexpetrean80/teleport-ui/internal/ui"
)

func main() {
	ctx := context.Background()

	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[0]
	filterArgs := args[1:]

	switch cmd {
	case "-h", "-help", "--help":
		printUsage()
		os.Exit(0)
	case "-c", "--clear-cache":
		if err := cache.ClearAll(); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "Cache cleared.")
		os.Exit(0)
	case "db":
		runDB(ctx, filterArgs)
	case "kube":
		runKube(ctx, filterArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func isHelpFlag(s string) bool {
	return s == "-h" || s == "-help" || s == "--help"
}

func extractClearCache(args []string) (bool, []string) {
	clearCache := false
	filtered := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--clear-cache" || a == "-c" {
			clearCache = true
		} else {
			filtered = append(filtered, a)
		}
	}
	return clearCache, filtered
}

func runDB(ctx context.Context, filterArgs []string) {
	if len(filterArgs) > 0 && isHelpFlag(filterArgs[0]) {
		printDBUsage()
		os.Exit(0)
	}

	clearCache, filterArgs := extractClearCache(filterArgs)

	dbs, err := teleport.GetTeleportDatabases(ctx, filterArgs, clearCache)
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

	dbName, confirmed, err := ui.RunTextInput("Database name", selectedDB.Metadata.Labels.DBName)
	if err != nil {
		log.Fatal(err)
	}
	if !confirmed {
		fmt.Println("Cancelled")
		return
	}

	if err = teleport.ConnectToTeleportDB(ctx, selectedDB, *selectedDBUser, dbName); err != nil {
		log.Fatal(err)
	}
}

func runKube(ctx context.Context, filterArgs []string) {
	if len(filterArgs) > 0 && isHelpFlag(filterArgs[0]) {
		printKubeUsage()
		os.Exit(0)
	}

	clearCache, filterArgs := extractClearCache(filterArgs)

	clusters, err := teleport.GetTeleportKubeClusters(ctx, filterArgs, clearCache)
	if err != nil {
		log.Fatal(err)
	}

	selectedCluster, err := ui.RunFuzzyFinder(clusters)
	if err != nil {
		log.Fatal(err)
	}

	if selectedCluster == nil {
		fmt.Println("No cluster selected")
		return
	}

	fmt.Printf("Selected: %s\n", selectedCluster.String())

	if err = teleport.LoginToTeleportKube(ctx, selectedCluster); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: teleport-ui <command> [flags] [filter args...]\n")
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  db    List and connect to databases\n")
	fmt.Fprintf(os.Stderr, "  kube  List and login to Kubernetes clusters\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  -c, --clear-cache  Clear all cached tsh data and exit\n")
	fmt.Fprintf(os.Stderr, "  -h, --help         Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nUse \"teleport-ui <command> --help\" for more information about a command.\n")
}

func printDBUsage() {
	fmt.Fprintf(os.Stderr, "Usage: teleport-ui db [flags] [filter args...]\n")
	fmt.Fprintf(os.Stderr, "\nList Teleport databases and connect to a selected one.\n")
	fmt.Fprintf(os.Stderr, "A fuzzy finder lets you pick the database and user.\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  -c, --clear-cache      Clear the database cache and fetch fresh data\n")
	fmt.Fprintf(os.Stderr, "      --search=<query>   Search for databases by name or label (passed to tsh)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help             Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nFilter args are passed to tsh (e.g. key1=value1, --search=foo, --query='...')\n")
	fmt.Fprintf(os.Stderr, "\nResults are cached in ~/.cache/teleport-ui/ for faster subsequent runs.\n")
	fmt.Fprintf(os.Stderr, "Use --clear-cache when new resources have been added.\n")
}

func printKubeUsage() {
	fmt.Fprintf(os.Stderr, "Usage: teleport-ui kube [flags] [filter args...]\n")
	fmt.Fprintf(os.Stderr, "\nList Teleport Kubernetes clusters and login to a selected one.\n")
	fmt.Fprintf(os.Stderr, "A fuzzy finder lets you pick the cluster.\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  -c, --clear-cache      Clear the kube cache and fetch fresh data\n")
	fmt.Fprintf(os.Stderr, "      --search=<query>   Search for clusters by name or label (passed to tsh)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help             Show this help message\n")
	fmt.Fprintf(os.Stderr, "\nFilter args are passed to tsh (e.g. key1=value1, --search=foo, --query='...')\n")
	fmt.Fprintf(os.Stderr, "\nResults are cached in ~/.cache/teleport-ui/ for faster subsequent runs.\n")
	fmt.Fprintf(os.Stderr, "Use --clear-cache when new resources have been added.\n")
}
