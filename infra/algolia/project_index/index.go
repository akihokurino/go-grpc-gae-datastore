package project_index

import (
	"gae-go-sample/domain"
)

const indexName = "project"

type index struct {
	ObjectID                  string `json:"objectID"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	CompanyName               string `json:"companyName"`
	CompanyRepresentativeName string `json:"companyRepresentativeName"`
	CompanyIntroduction       string `json:"companyIntroduction"`
	AccordingCompanyName      string `json:"accordingCompanyName"`
	AccordingCompanyAddress   string `json:"accordingCompanyAddress"`
	CreatedAt                 int64  `json:"createdAt"`
	Status                    int32  `json:"status"`

	Highlighted map[string]interface{} `json:"_highlightResult"`
}

func newIndex(project *domain.Project, company *domain.Company) index {
	return index{
		ObjectID:                  project.ID.String(),
		Name:                      project.Name,
		Description:               project.Description,
		CompanyName:               company.Name,
		CompanyRepresentativeName: company.RepresentativeName,
		CompanyIntroduction:       company.Introduction,
		AccordingCompanyName:      company.AccordingCompanyName,
		AccordingCompanyAddress:   company.AccordingCompanyAddress,
		CreatedAt:                 project.CreatedAt.UTC().Unix(),
		Status:                    int32(project.Status),
	}
}
