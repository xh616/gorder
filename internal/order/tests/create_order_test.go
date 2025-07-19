package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	sw "github.com/xh/gorder/internal/common/client/order"
	_ "github.com/xh/gorder/internal/common/config"
)

var (
	ctx    = context.Background()
	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
)

func TestMain(m *testing.M) {
	before()
	m.Run()
}

func before() {
	log.Printf("server=%s", server)
}

func TestCreateOrder_success(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S7EoKpdVdkkKYi",
				Quantity: 1,
			},
		},
	})
	t.Logf("body=%s", string(response.Body))
	assert.Equal(t, 200, response.StatusCode())

	assert.Equal(t, 0, response.JSON200.Errno)
}

func TestCreateOrder_invalidParams(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items:      nil,
	})
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 1, response.JSON200.Errno)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("getResponse body=%+v", body)
	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
