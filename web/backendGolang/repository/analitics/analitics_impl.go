package analiticsrepo

type AnaliticsRepository interface{
	GetAnalitics() (int64, error)
	GetAnaliticsPerDay() (map[string]int64, error)
}