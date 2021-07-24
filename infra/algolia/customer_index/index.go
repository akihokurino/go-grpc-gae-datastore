package customer_index

import (
	"gae-go-sample/domain"
)

const indexName = "customer"
const indexNameForClientOrderByCreatedAt = "customer_overview_order_by_created_at"
const indexNameForClientOrderByBookmark = "customer_overview_order_by_bookmark"

type index struct {
	ObjectID  string `json:"objectID"`
	Name      string `json:"name"`
	NameKana  string `json:"nameKana"`
	PR        string `json:"pr"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"createdAt"`
	Status    int32  `json:"status"`

	Highlighted map[string]interface{} `json:"_highlightResult"`
}

func newIndex(customer *domain.Customer) index {
	return index{
		ObjectID:  customer.ID.String(),
		Name:      customer.Name,
		NameKana:  customer.NameKana,
		PR:        customer.Pr,
		Address:   customer.Address,
		CreatedAt: customer.CreatedAt.UTC().Unix(),
		Status:    int32(customer.Status),
	}
}
