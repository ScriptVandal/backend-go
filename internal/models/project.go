package models

type Project struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Tags        []string `json:"tags"`
    URL         string   `json:"url"`
}
