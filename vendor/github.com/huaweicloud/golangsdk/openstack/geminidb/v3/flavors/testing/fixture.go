package testing

import (
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/flavors"
)

const AllInstancesFlavorResponse = `
{ 
    "total_count": 4,
    "flavors": [ 
       { 
            "engine_name": "GeminiDB-Cassandra", 
            "engine_version": "3.11",
            "vcpus": "4", 
            "ram": "16", 
            "spec_code": "nosql.cassandra.xlarge.4", 
            "availability_zone": [ 
                "az1", 
                "az2" 
            ] 
        }, 
       { 
            "engine_name": "GeminiDB-Cassandra", 
            "engine_version": "3.11",
            "vcpus": "8", 
            "ram": "32", 
            "spec_code": "nosql.cassandra.2xlarge.4", 
            "availability_zone": [ 
                "az1", 
                "az2" 
            ] 
        }, 
       { 
            "engine_name": "GeminiDB-Cassandra", 
            "engine_version": "3.11",
            "vcpus": "16", 
            "ram": "64", 
            "spec_code": "nosql.cassandra.4xlarge.4", 
            "availability_zone": [ 
                "az1", 
                "az2" 
            ] 
        },
       { 
            "engine_name": "GeminiDB-Cassandra", 
            "engine_version": "3.11",
            "vcpus": "32", 
            "ram": "128", 
            "spec_code": "nosql.cassandra.8xlarge.4", 
            "availability_zone": [ 
                "az1", 
                "az2" 
            ] 
        } 
    ] 
}
`

var ExpectedAllInstancesFlavorResponse = flavors.AllInstancesFlavorResponse{
	TotalCount: 4,
	Flavors: []flavors.Flavors{

		{
			EngineName:    "GeminiDB-Cassandra",
			EngineVersion: "3.11",
			Vcpus:         "4",
			Ram:           "16",
			SpecCode:      "nosql.cassandra.xlarge.4",
			AvailabilityZone: []string{
				"az1",
				"az2",
			},
		},

		{
			EngineName:    "GeminiDB-Cassandra",
			EngineVersion: "3.11",
			Vcpus:         "8",
			Ram:           "32",
			SpecCode:      "nosql.cassandra.2xlarge.4",
			AvailabilityZone: []string{
				"az1",
				"az2",
			},
		},

		{
			EngineName:    "GeminiDB-Cassandra",
			EngineVersion: "3.11",
			Vcpus:         "16",
			Ram:           "64",
			SpecCode:      "nosql.cassandra.4xlarge.4",
			AvailabilityZone: []string{
				"az1",
				"az2",
			},
		},
		{
			EngineName:    "GeminiDB-Cassandra",
			EngineVersion: "3.11",
			Vcpus:         "32",
			Ram:           "128",
			SpecCode:      "nosql.cassandra.8xlarge.4",
			AvailabilityZone: []string{
				"az1",
				"az2",
			},
		},
	},
}
