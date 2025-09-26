package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var SupaClient *supabase.Client

func InitDB() {
	if os.Getenv("VERCEL") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found (this is normal in production)")
		}
	}

	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_ANON_KEY")

	// DEBUG: Uncomment these lines to see what's being loaded
	// log.Printf("SUPABASE_URL: '%s'", supaURL)
	// log.Printf("SUPABASE_ANON_KEY length: %d", len(supaKey))

	if supaURL == "" || supaKey == "" {
		log.Fatal("Missing SUPABASE_URL or SUPABASE_ANON_KEY environment variables")
	}

	client, err := supabase.NewClient(supaURL, supaKey, nil)
	if err != nil {
		log.Fatal("Failed to initialize Supabase client:", err)
	}

	SupaClient = client
}
