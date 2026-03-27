export function DashboardPage() {
  return (
    <div>
      <div className="page-header">
        <div>
          <div className="eyebrow">Dashboard</div>
          <h1>每日分析总览</h1>
          <p className="muted">聚焦 A 股每日报告生产、任务执行和历史结果查看。</p>
        </div>
      </div>

      <div className="card-grid">
        <section className="panel">
          <h3>今日任务</h3>
          <p>每日批量分析任务将在这里展示执行状态、完成数量和失败重试情况。</p>
        </section>
        <section className="panel">
          <h3>最近报告</h3>
          <p>最近生成的 A 股每日报告会在这里展示摘要，便于快速复盘。</p>
        </section>
        <section className="panel">
          <h3>模型配置</h3>
          <p>当前版本仅支持 OpenAI 兼容 API，可配置 Base URL、模型名和 API Key。</p>
        </section>
      </div>
    </div>
  )
}
