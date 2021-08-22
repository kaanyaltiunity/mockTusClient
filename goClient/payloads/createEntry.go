package payloads

type CreateEntry struct {
	Path        string `json:"path"`
	ContentHash string `json:"conetnt_hash"`
	ContentSize string `json:"content_size"`
	ContentType string `json:"content_type"`
}
