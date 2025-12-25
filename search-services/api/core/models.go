package core

type UpdateStatus string

const (
	StatusUpdateUnknown UpdateStatus = "unknown"
	StatusUpdateIdle    UpdateStatus = "idle"
	StatusUpdateRunning UpdateStatus = "running"
)

type UpdateStats struct {
	WordsTotal    int
	WordsUnique   int
	ComicsFetched int
	ComicsTotal   int
}

type Comics struct {
	ID    int
	URL   string
	Words []string
}

type Comic struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type SearchResult struct {
	Comics []Comic `json:"comics"`
	Total  int64   `json:"total"`
}
