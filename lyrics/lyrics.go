package lyrics

import (
	"context"
	"log"

	language "cloud.google.com/go/language/apiv1"
)

// Malappropriate takes some lyrics, then malappropriates them by substituting in various random
// rhymings words. It returns the lyrics in the same format it was passed them.
func Malappropriate(lyrics string) string {
	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
}
