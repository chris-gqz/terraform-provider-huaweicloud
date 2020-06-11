package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
	th "github.com/huaweicloud/golangsdk/testhelper"
	"github.com/huaweicloud/golangsdk/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, CreateResponse)
	})

	options := instances.CreateGeminiDBOpts{
		Name: "test-cassandra-01",
		Datastore: instances.DataStore{
			Type:          "GeminiDB-Cassandra",
			Version:       "3.11",
			StorageEngine: "rocksDB",
		},
		Region:           "aaa",
		AvailabilityZone: "bbb",
		VpcId:            "674e9b42-cd8d-4d25-a2e6-5abcc565b961",
		SubnetId:         "f1df08c5-71d1-406a-aff0-de435a51007b",
		SecurityGroupId:  "7aa51dbf-5b63-40db-9724-dad3c4828b58",
		Password:         "Test@123",
		Mode:             "Cluster",
		Flavor: []instances.Flavor{
			{
				Num:      3,
				Size:     500,
				Storage:  "ULTRAHIGH",
				SpecCode: "nosql.cassandra.4xlarge.4",
			},
		},
		BackupStrategy: &instances.BackupStrategy{
			StartTime: "08:15-09:15",
			KeepDays:  8,
		},
		EnterpriseProjectId: "0",
	}

	actual, err := instances.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	expected := instances.CreateGeminiDB{
		AvailabilityZone: "bbb",
		Flavor: []instances.Flavor{
			{
				Num:      3,
				Size:     500,
				Storage:  "ULTRAHIGH",
				SpecCode: "nosql.cassandra.4xlarge.4",
			},
		},
		JobId: "c010abd0-48cf-4fa8-8cbc-090f093eaa2f",
	}

	expected.GeminiDBBase = instances.GeminiDBBase{
		Id:   "39b6a1a278844ac48119d86512e0000bin06",
		Name: "test-cassandra-01",
		Datastore: instances.DataStore{
			Type:          "GeminiDB-Cassandra",
			Version:       "3.11",
			StorageEngine: "rocksDB",
		},
		Created:             "2019-10-28 14:10:54",
		Status:              "creating",
		Region:              "aaa",
		VpcId:               "674e9b42-cd8d-4d25-a2e6-5abcc565b961",
		SubnetId:            "f1df08c5-71d1-406a-aff0-de435a51007b",
		SecurityGroupId:     "7aa51dbf-5b63-40db-9724-dad3c4828b58",
		Mode:                "Cluster",
		EnterpriseProjectId: "0",
		BackupStrategy: instances.BackupStrategy{
			StartTime: "08:15-09:15",
			KeepDays:  8,
		},
	}
	th.AssertDeepEquals(t, expected, *actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances/4e8e5957", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
	res := instances.Delete(client.ServiceClient(), "4e8e5957")
	th.AssertNoErr(t, res.Err)
}
