package service

import "context"

type IReport interface {
	List(ctx context.Context) ([]ReportItem, error)
}

type ReportItem struct {
	ID             uint   `json:"id"`
	TradeDate      string `json:"tradeDate"`
	StockCode      string `json:"stockCode"`
	StockName      string `json:"stockName"`
	Summary        string `json:"summary"`
	RiskLevel      string `json:"riskLevel"`
	Recommendation string `json:"recommendation"`
}

var localReport IReport

func Report() IReport {
	return localReport
}

func RegisterReport(s IReport) {
	localReport = s
}
