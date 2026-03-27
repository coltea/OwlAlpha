import { useEffect, useState } from 'react'
import { fetchReports, type ReportItem } from '../../services/report'

export function ReportsPage() {
  const [items, setItems] = useState<ReportItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchReports()
      .then(setItems)
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <div className="page-header">
        <div>
          <div className="eyebrow">Reports</div>
          <h1>报告中心</h1>
          <p className="muted">查看每日生成的 A 股分析报告。</p>
        </div>
      </div>

      <section className="panel">
        {loading ? <p>加载中...</p> : null}
        {!loading && items.length === 0 ? <p>暂无报告数据。</p> : null}
        <div className="report-list">
          {items.map((item) => (
            <article key={item.id} className="report-item">
              <div className="report-meta">
                <strong>{item.stockName}</strong>
                <span>{item.stockCode}</span>
                <span>{item.tradeDate}</span>
              </div>
              <p>{item.summary}</p>
              <div className="badge-row">
                <span className="badge">风险: {item.riskLevel}</span>
                <span className="badge">建议: {item.recommendation}</span>
              </div>
            </article>
          ))}
        </div>
      </section>
    </div>
  )
}
