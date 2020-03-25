package financial

import "time"

type Financial struct {
	Date                     time.Time `bson:"date" json:"date"`
	AvarageSharesOutstanding float64   `bson:"avarageSharesOutstanding" json:"avarageSharesOutstanding"`
	StockholdersEquity       float64   `bson:"stockholdersEquity" json:"stockholdersEquity"`
	TotalRevenue             float64   `bson:"totalRevenue" json:"totalRevenue"`
	NetIncome                float64   `bson:"netIncome" json:"netIncome"`
	EBIT                     float64   `bson:"ebit" json:"ebit"`
}
