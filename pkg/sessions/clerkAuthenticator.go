package sessions

import (
	"log"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type ClerkAuthenticator struct {
	Client clerk.Client
}

func NewClerkAuthenticator() *ClerkAuthenticator {
	client, err := clerk.NewClient(os.Getenv("CLERK_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	return &ClerkAuthenticator{
		Client: client,
	}
}
