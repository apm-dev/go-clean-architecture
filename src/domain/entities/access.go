package entities

type UserID int64

type Method string

const (
	CreateBlog Method = "blog.create"
)

func AllMethods() map[Method]bool {
	return map[Method]bool{
		CreateBlog: true,
	}
}
