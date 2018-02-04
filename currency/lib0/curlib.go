package curlib

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Currency struct {
	Code    string
	Name    string
	Number  string
	Country string
}

func Load(path string) []Currency {
	table := make([]Currency, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		c := Currency{
			Country: row[0],
			Name:    row[1],
			Code:    row[2],
			Number:  row[3],
		}
		table = append(table, c)
	}
	return table
}

func Find(table []Currency, filter string) []Currency {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]Currency, 0)
	filter = strings.ToUpper(filter)
	for _, cur := range table {
		if cur.Code == filter ||
			cur.Number == filter ||
			strings.Contains(strings.ToUpper(cur.Country), filter) ||
			strings.Contains(strings.ToUpper(cur.Name), filter) {
			result = append(result, cur)
		}
	}
	return result
}
