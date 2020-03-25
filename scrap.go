package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/brunofortes/b3-companies-income-statements/scraper"
)

func main() {
	fmt.Println("Scraping for companies income statement.")

	var period = Annual
	if os.Args[1] == "quarterly" {
		period = Quarterly
	}

	periodBegin, _ := time.Parse("2006-01-02", os.Args[2])
	periodEnd, _ := time.Parse("2006-01-02", os.Args[3])

	scraper := Scraper{}
	scraper.ScrapCompanies(period, periodBegin, periodEnd)

	fmt.Println("Finish.")
}
