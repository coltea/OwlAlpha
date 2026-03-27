import { http } from './http'

export type ReportItem = {
  id: number
  tradeDate: string
  stockCode: string
  stockName: string
  summary: string
  riskLevel: string
  recommendation: string
}

export async function fetchReports(): Promise<ReportItem[]> {
  const response = await http.get('/reports')
  return response.data.data.items
}
