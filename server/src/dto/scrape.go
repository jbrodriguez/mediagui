package dto

// Scrape -
type Scrape struct {
	BaseURL       string
	SecureBaseURL string
	// BasePath      string
	Movie  interface{}
	Forced bool
}
