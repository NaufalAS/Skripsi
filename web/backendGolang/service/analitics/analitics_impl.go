package analiticsservice

type AnaliticsService interface{
GetAnalyticsData() (map[string]interface{}, error)
}