package services

import "github.com/adwinugroho/go-rest-self-screening/models"

// Create a interface to implements our function at models
type (
	impelementsModels interface {
		AddData(data models.HealthAssessment) (*models.HealthAssessment, error)
		DeleteByKey(id string) (*[]models.HealthAssessment, error)
		GetDataByKey(id string) (*models.HealthAssessment, error)
		GetListAllData(vars models.BindVars) ([]models.HealthAssessment, int64, error)
		SaveLog(model *models.HealthAssessment) (*string, error)
		UpdateData(model *models.HealthAssessment) (*models.HealthAssessment, error)
	}
	// struct for reciever
	HealthService struct {
		dao impelementsModels
	}
)

// Create a function to call implements models first
func NewService(dao impelementsModels) *HealthService {
	return &HealthService{dao}
}

func (service *HealthService) SubmitCovidScreening(reqData models.HealthAssessment) (*models.HealthAssessment, error) {
	return &reqData, nil
}
