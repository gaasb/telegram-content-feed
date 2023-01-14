package bot

type Model struct {
	Categorys Category
	UserID    uint64
	PostID    uint64
	VideosURL map[string]string
	ExpireAt  *string
}
