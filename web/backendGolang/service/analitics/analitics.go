package analiticsservice

import analiticsrepo "skripsi/repository/analitics"

type AnaliticsServiceImpl struct {
	analiticsrrepo analiticsrepo.AnaliticsRepository
}

func NewAnaliticsService(analiticsrrepo analiticsrepo.AnaliticsRepository) *AnaliticsServiceImpl {
	return &AnaliticsServiceImpl{
		analiticsrrepo: analiticsrrepo,
	}
}

func (service *AnaliticsServiceImpl) GetAnalyticsData() (map[string]interface{}, error) {
	// Map untuk menyimpan hasil total data dan data per hari
	result := make(map[string]interface{})

	// Panggil fungsi GetAnalitics dari repo untuk menghitung total data
	analyticsTotal, err := service.analiticsrrepo.GetAnalitics()
	if err != nil {
		return nil, err
	}

	// Menambahkan total data ke map
	result["total_data"] = analyticsTotal

	// Panggil fungsi GetAnaliticsPerDay dari repo untuk menghitung data per hari
	analyticsPerDay, err := service.analiticsrrepo.GetAnaliticsPerDay()
	if err != nil {
		return nil, err
	}

	// Menambahkan data per hari ke map
	result["data_per_day"] = analyticsPerDay

	// Mengembalikan hasil dalam bentuk map yang berisi total data dan data per hari
	return result, nil
}
