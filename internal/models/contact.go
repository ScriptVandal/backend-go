package models

type Contact struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Telegram string `json:"telegram"`
    LinkedIn string `json:"linkedin"`
    Github   string `json:"github"`
}
