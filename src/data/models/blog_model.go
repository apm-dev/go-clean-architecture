package models

type BlogModel struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	CreatedAt int64
	UpdateAt  int64
}
