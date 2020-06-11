package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*CreateGeminiDB, error) {
	var response CreateGeminiDB
	err := r.ExtractInto(&response)
	return &response, err
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*ListGeminiDBResponse, error) {
	var response ListGeminiDBResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type DeleteInstanceGeminiDBResult struct {
	commonResult
}

type DeleteInstanceGeminiDBResponse struct {
	JobId string `json:"job_id"`
}

func (r DeleteInstanceGeminiDBResult) Extract() (*DeleteInstanceGeminiDBResponse, error) {
	var response DeleteInstanceGeminiDBResponse
	err := r.ExtractInto(&response)
	return &response, err
}

type ListGeminiDBResult struct {
	commonResult
}

type ListGeminiDBResponse struct {
	Instances  []GeminiDBInstanceResponse `json:"instances"`
	TotalCount int                        `json:"total_count"`
}

type GeminiDBBase struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Region          string    `json:"region"`
	Mode            string    `json:"mode"`
	DataStore       DataStore `json:"datastore"`
	Created         string    `json:"created"`
	VpcId           string    `json:"vpc_id"`
	SubnetId        string    `json:"subnet_id"`
	SecurityGroupId string    `json:"security_group_id"`

	EnterpriseProjectId string `json:"enterprise_project_id"`
	AvailabilityZone    string `json:"availability_zone"`
}

type BackupStrategyList struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  int    `json:"keep_days,omitempty"`
}

type CreateGeminiDB struct {
	Flavor         []Flavor       `json:"flavor"`
	JobId          string         `json:"job_id"`
	BackupStrategy BackupStrategy `json:"backup_strategy"`
	GeminiDBBase
}

type GeminiDBInstanceResponse struct {
	GeminiDBBase
	Port              string             `json:"port"`
	Engine            string             `json:"engine"`
	Updated           string             `json:"updated"`
	DbUserName        string             `json:"db_user_name"`
	PayMode           string             `json:"pay_mode"`
	MaintenanceWindow string             `json:"maintenance_window"`
	Groups            []Groups           `json:"groups"`
	TimeZone          string             `json:"time_zone"`
	Actions           []string           `json:"actions"`
	BackupStrategy    BackupStrategyList `json:"backup_strategy"`
}

type Groups struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Volume Volume  `json:"volume"`
	Nodes  []Nodes `json:"nodes"`
}

type Volume struct {
	Size string `json:"size"`
	used string `json:"used"`
}

type Nodes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	PrivateIp        string `json:"private_ip"`
	SpecCode         string `json:"spec_code"`
	AvailabilityZone string `json:"availability_zone"`
}

type GeminiDBPage struct {
	pagination.SinglePageBase
}

func (r GeminiDBPage) IsEmpty() (bool, error) {
	data, err := ExtractGeminiDBInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractGeminiDBInstances is a function that takes a ListResult and returns the services' information.
func ExtractGeminiDBInstances(r pagination.Page) (ListGeminiDBResponse, error) {
	var s ListGeminiDBResponse
	err := (r.(GeminiDBPage)).ExtractInto(&s)
	return s, err
}
