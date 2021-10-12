package payloads

type CreateEntry struct {
	Path        string `json:"path"`
	ContentHash string `json:"content_hash"`
	ContentSize int    `json:"content_size"`
	ContentType string `json:"content_type"`
}
