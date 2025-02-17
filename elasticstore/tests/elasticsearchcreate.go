package tests

import (
	"fmt"

	esClient "github.com/kiranpt03/factorysphere/devicesphere/engines/elasticstore/pkg/databases/elasticsearch"
)

type NestedData struct {
	ID        string `json:"id"`
	Timestamp string `json:"@timestamp"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Address   struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		Zipcode string `json:"zipcode"`
	} `json:"address"`
	Interests []string `json:"interests"`
}

func main() {
	client := esClient.GetClient("http://localhost:9200")

	data := NestedData{
		ID:        "1",
		Timestamp: "",
		Name:      "John Doe",
		Age:       30,
		Address: struct {
			Street  string `json:"street"`
			City    string `json:"city"`
			State   string `json:"state"`
			Zipcode string `json:"zipcode"`
		}{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			Zipcode: "12345",
		},
		Interests: []string{"reading", "hiking", "coding"},
	}

	res, err := client.Create("test", data.ID, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
