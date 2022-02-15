package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/spf13/cobra"
)

const version = "4.3.0"

var rootCmd = &cobra.Command{
	Use:     "openbooks",
	Short:   "Quickly and easily download eBooks from IRCHighway.",
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Generate a random username to avoid IRC name collisions if multiple users are hosting
// at the same time.
func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(int64(rand.Int()))
	return fmt.Sprintf("%s_%s", gofakeit.Adjective(), gofakeit.Noun())
}
