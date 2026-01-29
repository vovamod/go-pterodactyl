package main

import (
	"context"
	"fmt"
	"github.com/vovamod/go-pterodactyl"
	"log"
	"os"
	"time"
)

func main() {

	baseURL := os.Getenv("PTERO_BASE_URL")
	apiKey := os.Getenv("PTERO_API_KEY")
	if baseURL == "" || apiKey == "" {
		log.Fatal("Set PTERO_BASE_URL and PTERO_API_KEY environment variables")
	}

	// Initialise SDK with a short default timeout; callers can override via ctx.
	client, err := pterodactyl.NewClient(baseURL, apiKey, pterodactyl.ApplicationKey, pterodactyl.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	// Always pass a context so the request can be cancelled.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	nodes, err := client.ApplicationAPI.Nodes.ListAll(ctx)
	if err != nil {
		log.Fatalf("list nodes: %v", err)
	}

	fmt.Printf("Found %d nodes:\n", len(nodes))
	for _, n := range nodes {
		fmt.Printf("• ID %d – %-20s Memory: %dMB  Disk: %dMB\n", n.ID, n.Name, n.Memory, n.Disk)
	}

}
