package core

import (
	"sort"
	"sync"
)

type Index struct {
	mu    sync.RWMutex
	items map[string][]int64 // keyword -> []comicID
	docs  map[int64]Comic    // comicID -> Comic
}

func NewIndex() *Index {
	return &Index{
		items: make(map[string][]int64),
		docs:  make(map[int64]Comic),
	}
}

func (i *Index) Add(comics []Comic) {
	i.mu.Lock()
	defer i.mu.Unlock()

	newItems := make(map[string][]int64)
	newDocs := make(map[int64]Comic)

	for _, comic := range comics {
		newDocs[comic.ID] = comic
		for _, keyword := range comic.Keywords {
			newItems[keyword] = append(newItems[keyword], comic.ID)
		}
	}

	for k := range newItems {
		sort.Slice(newItems[k], func(a, b int) bool {
			return newItems[k][a] < newItems[k][b]
		})
	}

	i.items = newItems
	i.docs = newDocs
}

func (i *Index) Search(keywords []string) []Comic {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if len(keywords) == 0 {
		return nil
	}

	matches := make(map[int64]int)

	for _, kw := range keywords {
		if ids, ok := i.items[kw]; ok {
			for _, id := range ids {
				matches[id]++
			}
		}
	}

	type match struct {
		ID    int64
		Count int
	}
	var result []match
	for id, count := range matches {
		result = append(result, match{ID: id, Count: count})
	}

	sort.Slice(result, func(a, b int) bool {
		if result[a].Count != result[b].Count {
			return result[a].Count > result[b].Count
		}
		return result[a].ID < result[b].ID
	})

	comics := make([]Comic, len(result))
	for k, v := range result {
		comics[k] = i.docs[v.ID]
	}

	return comics
}
