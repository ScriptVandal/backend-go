package models

type Post struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Content     string   `json:"content"`
    Tags        []string `json:"tags"`
    PublishedAt string   `json:"published_at"`
}
