package core

type Comic struct {
	ID       int64
	URL      string
	Keywords []string
}

type SearchResult struct {
	Comics []Comic
	Total  int64
}
