package report

import (
	"context"

	"github.com/coltea/owlalpha/backend/internal/bootstrap"
	"github.com/coltea/owlalpha/backend/internal/model/entity"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type Logic struct {
	deps *bootstrap.Dependencies
}

func New(deps *bootstrap.Dependencies) *Logic {
	return &Logic{deps: deps}
}

func (l *Logic) List(ctx context.Context) ([]service.ReportItem, error) {
	var reports []entity.Report
	if err := l.deps.DB.WithContext(ctx).Order("trade_date desc, id desc").Limit(50).Find(&reports).Error; err != nil {
		return nil, err
	}

	items := make([]service.ReportItem, 0, len(reports))
	for _, item := range reports {
		items = append(items, service.ReportItem{
			ID:             item.ID,
			TradeDate:      item.TradeDate,
			StockCode:      item.StockCode,
			StockName:      item.StockName,
			Summary:        item.Summary,
			RiskLevel:      item.RiskLevel,
			Recommendation: item.Recommendation,
		})
	}

	return items, nil
}
