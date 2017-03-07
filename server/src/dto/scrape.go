package dto

type Scrape struct {
	BaseUrl       string
	SecureBaseUrl string
	// BasePath      string
	Movie  interface{}
	Forced bool
}
