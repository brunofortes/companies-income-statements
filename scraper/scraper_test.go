package scraper

import (
	"testing"
	"time"

	. "github.com/brunofortes/b3-companies-income-statements/financial"
)

func TestScraper(t *testing.T) {
	assertEqualStrings := func(t *testing.T, got string, expected string, fieldName string) {
		t.Helper()

		if got != expected {
			t.Errorf("got '%s' expected '%s' as %s.", got, expected, fieldName)
		}
	}

	containsElement := func(slice []string, find string) bool {
		for _, f := range slice {
			if f == find {
				return true
			}
		}

		return false
	}

	assertContainsElement := func(t *testing.T, got []string, expected string, fieldName string) {
		t.Helper()

		if got == nil || len(got) == 0 {
			t.Errorf("got an empty slice for %s.", fieldName)
		}

		if !containsElement(got, expected) {
			t.Errorf("got '%s' expected a slice containing '%s' as %s.", got, expected, fieldName)
		}
	}

	assertEqualFloat := func(t *testing.T, got float64, expected float64, fieldName string) {
		t.Helper()

		if got != expected {
			t.Errorf("got '%f' expected '%f' as %s.", got, expected, fieldName)
		}
	}

	scraper := Scraper{}
	t.Run("must obtain company profile information", func(t *testing.T) {
		expectedName := "Petroleo Brasileiro S.A. - Petrobras"
		expectedSector := "Energy"
		expectedIndustry := "Oil & Gas Integrated"
		expectedLabel := "PETR4"
		expectedLabel2 := "PETR"

		companyResult, err := scraper.ScrapCompanyProfile("PETR4.SA")
		if err != nil {
			t.Errorf("error getting company profile info")
		}

		assertEqualStrings(t, companyResult.Name, expectedName, "company name")
		assertContainsElement(t, companyResult.Sectors, expectedSector, "company sector")
		assertContainsElement(t, companyResult.Industries, expectedIndustry, "company industry")
		assertContainsElement(t, companyResult.Labels, expectedLabel, "company label")
		assertContainsElement(t, companyResult.Labels, expectedLabel2, "company label")
	})

	t.Run("must return ErrorCompanyNotFound when invalid symbol", func(t *testing.T) {
		var symbol = "XYZ"
		_, err := scraper.ScrapCompanyProfile(symbol)
		if err != ErrorCompanyNotFound {
			t.Errorf("expected ErrorCompanyNotFound for symbol %s", symbol)
		}
	})

	t.Run("must return ErrorCompanyNotFound when no industry and sector for symbol", func(t *testing.T) {
		var symbol = "ELET11.SA"
		_, err := scraper.ScrapCompanyProfile(symbol)
		if err != ErrorCompanyNotFound {
			t.Errorf("expected ErrorCompanyNotFound for symbol %s", symbol)
		}
	})

	findFinancialByDate := func(time time.Time, financials []Financial) Financial {
		for _, f := range financials {
			if f.Date == time {
				return f
			}
		}

		return Financial{}
	}

	t.Run("must obtain company financial information - annual", func(t *testing.T) {
		symbol := "UGPA3.SA"
		periodBegin, _ := time.Parse("2006-01-02", "2016-12-29")
		periodEnd, _ := time.Parse("2006-01-02", "2017-01-01")

		expectedDate, _ := time.Parse("2006-01-02", "2016-12-31")
		expectedAvarageSharesOutstanding := 1112810192.0
		expectedStockholdersEquity := 8527623000.0
		expectedTotalRevenue := 77352955000.0
		expectedNetIncome := 1561585000.0
		expectedEBIT := 4533536000.0

		financials, err := scraper.ScrapCompanyFinancial(symbol, Annual, periodBegin, periodEnd)
		if err != nil {
			t.Errorf("error getting company financial info")
		}

		financial2016 := findFinancialByDate(expectedDate, financials)
		assertEqualFloat(t, financial2016.AvarageSharesOutstanding, expectedAvarageSharesOutstanding, "AvarageSharesOutstanding")
		assertEqualFloat(t, financial2016.StockholdersEquity, expectedStockholdersEquity, "StockholdersEquity")
		assertEqualFloat(t, financial2016.TotalRevenue, expectedTotalRevenue, "TotalRevenue")
		assertEqualFloat(t, financial2016.NetIncome, expectedNetIncome, "NetIncome")
		assertEqualFloat(t, financial2016.EBIT, expectedEBIT, "EBIT")
	})

	t.Run("must return empty financial information when no financial for symbol", func(t *testing.T) {
		var symbol = "FNAM11.SA"
		periodBegin, _ := time.Parse("2006-01-02", "2016-12-29")
		periodEnd, _ := time.Parse("2006-01-02", "2017-01-01")

		financials, err := scraper.ScrapCompanyFinancial(symbol, Annual, periodBegin, periodEnd)
		if err != ErrorCompanyNotFound && len(financials) > 0 {
			t.Errorf("expected ErrorCompanyNotFound and no financial info for symbol '%s'.", symbol)
		}
	})

	t.Run("must obtain company financial information - quarterly", func(t *testing.T) {
		var symbol = "UGPA3.SA"
		periodBegin, _ := time.Parse("2006-01-02", "2018-12-20")
		periodEnd, _ := time.Parse("2006-01-02", "2019-06-30")

		expectedDate, _ := time.Parse("2006-01-02", "2018-12-31")
		expectedAvarageSharesOutstanding := 2116024000.0
		expectedStockholdersEquity := 9448105000.0
		expectedTotalRevenue := 23467044000.0
		expectedNetIncome := 507643000.0
		expectedEBIT := 1218622000.0

		financials, err := scraper.ScrapCompanyFinancial(symbol, Quarterly, periodBegin, periodEnd)
		if err != nil {
			t.Errorf("error getting company financial info")
		}

		financial2016 := findFinancialByDate(expectedDate, financials)
		assertEqualFloat(t, financial2016.AvarageSharesOutstanding, expectedAvarageSharesOutstanding, "AvarageSharesOutstanding")
		assertEqualFloat(t, financial2016.StockholdersEquity, expectedStockholdersEquity, "StockholdersEquity")
		assertEqualFloat(t, financial2016.TotalRevenue, expectedTotalRevenue, "TotalRevenue")
		assertEqualFloat(t, financial2016.NetIncome, expectedNetIncome, "NetIncome")
		assertEqualFloat(t, financial2016.EBIT, expectedEBIT, "EBIT")
	})

}
