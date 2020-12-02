package entities

type Blog struct {
	ID        string
	AuthorID  string
	Title     string
	Content   string
	CreatedAt int64
	UpdatedAt int64
}
