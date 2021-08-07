package prefecture

import "gae-go-recruiting-server/domain"

type root struct {
	Items []entity `json:"result"`
}

type entity struct {
	Code int32  `json:"prefCode"`
	Name string `json:"prefName"`
}

func (e *entity) toDomain() *domain.Prefecture {
	return &domain.Prefecture{Code: e.Code, Name: e.Name}
}
