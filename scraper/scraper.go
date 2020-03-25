package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	. "github.com/brunofortes/b3-companies-income-statements/company"
	"github.com/brunofortes/b3-companies-income-statements/financial"
	. "github.com/brunofortes/b3-companies-income-statements/financial"
	"github.com/gocolly/colly"
)

const (
	URL = "https://br.financas.yahoo.com/noticias/acoes-mais-negociadas" // B3
	//URL = "https://finance.yahoo.com/sector/ms_technology" // Technology
	//URL = https://finance.yahoo.com/sector/ms_energy // Energy
)

type Scraper struct {
}

func (s *Scraper) ScrapCompanies(period Period, periodBegin time.Time, periodEnd time.Time) {
	companyConnection := NewCompanyConnection()

	offset := 0
	i := s.visitAnotherPage(offset, 100, period, periodBegin, periodEnd, companyConnection)
	for i > 0 {
		offset += 100
		i = s.visitAnotherPage(offset, 100, period, periodBegin, periodEnd, companyConnection)
	}

	companyConnection.Disconnect()
}

func (s *Scraper) visitAnotherPage(offset int, count int, period Period, periodBegin time.Time, periodEnd time.Time, companyConnection CompanyConnection) int {
	c := colly.NewCollector()

	i := 0
	c.OnHTML("#scr-res-table tr", func(e *colly.HTMLElement) {
		e.ForEach("a", func(arg1 int, a *colly.HTMLElement) {
			i++
			company, err := s.ScrapCompanyProfile(a.Text)

			if err == nil {
				finds, _ := companyConnection.FindByName(company.Name)
				if len(finds) > 0 {
					company = finds[0]
					company.Labels = includeLabel(company.Labels, formatCompanySymbol(a.Text)[0])
				}

				financials, err := s.ScrapCompanyFinancial(a.Text, period, periodBegin, periodEnd)
				if err == nil {
					company.Financials = mergeFinancialInfos(company.Financials, financials)
					fmt.Println(a.Text, " - ", company.Name)
					err = companyConnection.Save(company)
					if err != nil {
						fmt.Println("-> Error for company:", company.Name)
						fmt.Println("->", err)
					}
				}
			}
		})
	})

	c.Visit(URL + "?offset=" + strconv.Itoa(offset) + "&count=" + strconv.Itoa(count))

	return i
}

var ErrorCompanyNotFound = errors.New("type")

func formatCompanySymbol(symbol string) []string {
	var result []string
	s := symbol

	if strings.Contains(symbol, ".SA") {
		s = strings.Split(symbol, ".SA")[0]
		result = append(result, s)
	}

	if len(s) > 4 {
		s = string([]rune(s)[:4])
	}

	result = append(result, s)
	return result
}

func mergeFinancialInfos(financials1 []Financial, financials2 []Financial) []Financial {
	result := []Financial{}
	mapp := map[time.Time]Financial{}

	for _, f := range financials1 {
		result = append(result, f)
		mapp[f.Date] = f
	}

	for _, f := range financials2 {
		_, present := mapp[f.Date]
		if !present {
			result = append(result, f)
		}
	}

	return result
}

func includeLabel(labels []string, label string) []string {
	for _, l := range labels {
		if l == label {
			return labels
		}
	}

	return append(labels, label)
}

func (s *Scraper) ScrapCompanyProfile(symbol string) (Company, error) {
	var company = Company{}

	res, err := http.Get("https://query2.finance.yahoo.com/v7/finance/quote?symbols=" + symbol)
	if err == nil {
		// fetch company full name
		var quoteResult map[string]map[string][]map[string]string
		json.NewDecoder(res.Body).Decode(&quoteResult)
		if res.StatusCode == 404 || len(quoteResult["quoteResponse"]["result"]) == 0 {
			return Company{}, ErrorCompanyNotFound
		}
		company.Name = quoteResult["quoteResponse"]["result"][0]["longName"]

		// fetch company industry and sector
		res, err = http.Get("https://query2.finance.yahoo.com/v10/finance/quoteSummary/" + symbol + "?modules=assetProfile")
		if err == nil {
			if res.StatusCode == 404 {
				return Company{}, ErrorCompanyNotFound
			}
			var quoteSummary map[string]map[string][]map[string]map[string]string
			json.NewDecoder(res.Body).Decode(&quoteSummary)
			company.Labels = append(company.Labels, formatCompanySymbol(symbol)...)

			if quoteSummary["quoteSummary"]["result"] != nil || len(quoteSummary["quoteSummary"]["result"]) > 0 {
				company.Sectors = append(company.Sectors, quoteSummary["quoteSummary"]["result"][0]["assetProfile"]["sector"])
				company.Industries = append(company.Industries, quoteSummary["quoteSummary"]["result"][0]["assetProfile"]["industry"])
			}
		}
	}

	return company, err
}

func (s *Scraper) ScrapCompanyFinancial(symbol string, period Period, periodBegin time.Time, periodEnd time.Time) ([]Financial, error) {
	var financials = []Financial{}

	res, err := http.Get("https://query2.finance.yahoo.com/ws/fundamentals-timeseries/v1/finance/timeseries/" + symbol + "?symbol=" + symbol + "&period1=" + strconv.FormatInt(periodBegin.Unix(), 10) + "&period2=" + strconv.FormatInt(periodEnd.Unix(), 10) + "&type=" + period.typesStringList)
	if err == nil {
		var financialResult map[string]map[string][]map[string]interface{}
		json.NewDecoder(res.Body).Decode(&financialResult)
		items := financialResult["timeseries"]["result"]
		if res.StatusCode == 404 || len(items) == 0 || items[0]["timestamp"] == nil {
			return financials, ErrorCompanyNotFound
		}

		for _, f := range items {
			typee := (f["meta"].(map[string]interface{})["type"]).([]interface{})[0].(string)
			if f[typee] != nil {
				dates := f[typee].([]interface{})

				for i, d := range dates {
					if len(financials) == 0 || len(financials) < i+1 {
						financials = append(financials, financial.Financial{})
					}
					financials[i].Date, _ = time.Parse("2006-01-02", d.(map[string]interface{})["asOfDate"].(string))
					financials[i] = period.typesMap[typee](financials[i], d.(map[string]interface{})["reportedValue"].(map[string]interface{})["raw"])
				}
			}
		}
	}

	return financials, err
}

type funcSetFinancialValue = func(financial Financial, value interface{}) Financial

type Period struct {
	typesMap        map[string]funcSetFinancialValue
	typesStringList string
}

var Annual = Period{typesMap: annualTypesMap, typesStringList: mapFinancialsTypesToStringList(annualTypesMap)}
var Quarterly = Period{typesMap: quarterlyTypesMap, typesStringList: mapFinancialsTypesToStringList(quarterlyTypesMap)}

func mapFinancialsTypesToStringList(types map[string]funcSetFinancialValue) string {
	keys := reflect.ValueOf(types).MapKeys()
	result := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		result[i] = keys[i].String()
	}

	return strings.Join(result, ",")
}

var annualTypesMap = map[string]funcSetFinancialValue{
	"annualBasicAverageShares": func(financial Financial, value interface{}) Financial {
		financial.AvarageSharesOutstanding = value.(float64)
		return financial
	},

	"annualStockholdersEquity": func(financial Financial, value interface{}) Financial {
		financial.StockholdersEquity = value.(float64)
		return financial
	},

	"annualTotalRevenue": func(financial Financial, value interface{}) Financial {
		financial.TotalRevenue = value.(float64)
		return financial
	},

	"annualNetIncome": func(financial Financial, value interface{}) Financial {
		financial.NetIncome = value.(float64)
		return financial
	},

	"annualEbitda": func(financial Financial, value interface{}) Financial {
		financial.EBIT = value.(float64)
		return financial
	},
}

var quarterlyTypesMap = map[string]funcSetFinancialValue{
	"quarterlyBasicAverageShares": func(financial Financial, value interface{}) Financial {
		financial.AvarageSharesOutstanding = value.(float64)
		return financial
	},

	"quarterlyStockholdersEquity": func(financial Financial, value interface{}) Financial {
		financial.StockholdersEquity = value.(float64)
		return financial
	},

	"quarterlyTotalRevenue": func(financial Financial, value interface{}) Financial {
		financial.TotalRevenue = value.(float64)
		return financial
	},

	"quarterlyNetIncome": func(financial Financial, value interface{}) Financial {
		financial.NetIncome = value.(float64)
		return financial
	},

	"quarterlyEbitda": func(financial Financial, value interface{}) Financial {
		financial.EBIT = value.(float64)
		return financial
	},
}
