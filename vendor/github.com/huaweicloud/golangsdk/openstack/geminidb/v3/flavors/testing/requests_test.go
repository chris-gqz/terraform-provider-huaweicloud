package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/flavors"
	"github.com/huaweicloud/golangsdk/pagination"
	"github.com/huaweicloud/golangsdk/testhelper/client"
	th "github.com/huaweicloud/golangsdk/testhelper"
)

func TestGetAllInstancesFlavor(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/flavors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AllInstancesFlavorResponse)
	})

	options := flavors.AllInstancesFlavor{
		EngineName: "GeminiDB-Cassandr",
		Region:     "aaa",
	}

	count := 0
	flavors.GetAllInstancesFlavor(client.ServiceClient(), options).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := flavors.ExtractAllInstancesFlavor(page)
		if err != nil {
			t.Errorf("Failed to extract instances: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedAllInstancesFlavorResponse, actual)

		return true, nil
	})
	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
