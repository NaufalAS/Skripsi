package helper

import (
	"net/url"
	"strconv"
	// "fmt"
)

func ExtractFilter(params url.Values) (string, int, int) {
	// Ambil filter tunggal dari URL
	filters := params.Get("search")
	// status, _:= strconv.Atoi(params.Get("status"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	page, _ := strconv.Atoi(params.Get("page"))

	// fmt.Println("Seacrh:", filters)
	return filters, limit, page
}
