package masterdata

type ItemContainer struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
