package domain

import "fmt"

type SortMode int

const (
	SortByFrequencyDesc SortMode = iota
	SortByNameAsc
)

type NameCount struct {
	Name  string
	Count int
}

type Config struct {
	FilePath string
	SortMode SortMode
}

type ResultSet []NameCount

func (r ResultSet) String() string {
	var result string
	for _, nc := range r {
		result += nc.String() + "\n"
	}
	return result
}

func (nc NameCount) String() string {
	//return nc.Name + ": " + string(rune(nc.Count+'0'))
	return fmt.Sprintf("%s: %d", nc.Name, nc.Count)
}
