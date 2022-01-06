package models

// Информация о форуме.
type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int32  `json:"posts,omitempty"`
	Threads int32  `json:"threads,omitempty"`
}
