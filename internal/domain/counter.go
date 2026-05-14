package domain

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

type Counter struct{}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Count(r io.Reader, mode SortMode) (ResultSet, error) {
	scanner := bufio.NewScanner(r)

	const maxLine = 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxLine)

	// собираем статистику
	counts := make(map[string]int)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}
		counts[name]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// конвертируем карту в слайс
	result := make(ResultSet, 0, len(counts))
	for name, count := range counts {
		result = append(result, NameCount{Name: name, Count: count})
	}

	c.sortResults(result, mode)

	return result, nil
}

func (c *Counter) sortResults(results ResultSet, mode SortMode) {
	switch mode {
	case SortByFrequencyDesc:
		sort.Slice(results, func(i, j int) bool {
			if results[i].Count != results[j].Count {
				return results[i].Count > results[j].Count
			}
			return results[i].Name < results[j].Name
		})
	case SortByNameAsc:
		sort.Slice(results, func(i, j int) bool {
			return results[i].Name < results[j].Name
		})
	default:
		// SortByFrequencyDesc по дефолту
		sort.Slice(results, func(i, j int) bool {
			if results[i].Count != results[j].Count {
				return results[i].Count > results[j].Count
			}
			return results[i].Name < results[j].Name
		})
	}
}
