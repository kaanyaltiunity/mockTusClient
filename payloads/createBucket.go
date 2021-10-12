package payloads

type CreateBucket struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ProjectGuid string `json:"projectguid"`
}
