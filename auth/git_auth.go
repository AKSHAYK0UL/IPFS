package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func GitAuth() (*github.Client, error) {
	godotenv.Load() // only used in the local environment
	accessToken := os.Getenv("ACCESSTOKEN")
	ctx := context.Background()
	fmt.Println(ctx)
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)
	return github.NewClient(httpClient), nil
}
