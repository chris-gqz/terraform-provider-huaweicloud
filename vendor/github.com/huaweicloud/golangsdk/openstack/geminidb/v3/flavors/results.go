package flavors

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type AllInstancesFlavorResult struct {
	commonResult
}

type AllInstancesFlavorResponse struct {
	Flavors    []Flavors `json:"flavors"`
	TotalCount int       `json:"total_count"`
}

type Flavors struct {
	EngineName       string   `json:"engine_name"`
	EngineVersion    string   `json:"engine_version"`
	Vcpus            string   `json:"vcpus"`
	Ram              string   `json:"ram"`
	SpecCode         string   `json:"spec_code"`
	AvailabilityZone []string `json:"availability_zone"`
}

type AllInstancesFlavorPage struct {
	pagination.SinglePageBase
}

func (r AllInstancesFlavorPage) IsEmpty() (bool, error) {
	data, err := ExtractAllInstancesFlavor(r)
	if err != nil {
		return false, err
	}
	return len(data.Flavors) == 0, err
}

// ExtractAllInstancesFlavor is a function that takes a ListResult and returns the services' information.
func ExtractAllInstancesFlavor(r pagination.Page) (AllInstancesFlavorResponse, error) {
	var s AllInstancesFlavorResponse
	err := (r.(AllInstancesFlavorPage)).ExtractInto(&s)
	return s, err
}
