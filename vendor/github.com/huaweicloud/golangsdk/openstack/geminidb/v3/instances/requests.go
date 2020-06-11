package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateGeminiDBOpts struct {
	Name                string         `json:"name"  required:"true"`
	DataStore           DataStore      `json:"datastore" required:"true"`
	Region              string         `json:"region" required:"true"`
	AvailabilityZone    string         `json:"availability_zone" required:"true"`
	VpcId               string         `json:"vpc_id" required:"true"`
	SubnetId            string         `json:"subnet_id" required:"true"`
	SecurityGroupId     string         `json:"security_group_id" required:"true"`
	Password            string         `json:"password" required:"true"`
	Mode                string         `json:"mode" required:"true"`
	Flavor              []Flavor       `json:"flavor" required:"true"`
	BackupStrategy      BackupStrategy `json:"backup_strategy,omitempty"`
	EnterpriseProjectId string         `json:"enterprise_project_id,omitempty"`
}

type DataStore struct {
	Type          string `json:"type" required:"true"`
	Version       string `json:"version" required:"true"`
	StorageEngine string `json:"storage_engine" required:"true"`
}

type Flavor struct {
	Num      string `json:"num" required:"true"`
	Size     string `json:"size" required:"true"`
	Storage  string `json:"storage" required:"true"`
	SpecCode string `json:"spec_code" required:"true"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  string `json:"keep_days,omitempty"`
}

type CreateGeminiDBBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateGeminiDBOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateGeminiDBBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteInstanceGeminiDBResult) {
	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes: []int{201, 202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})

	return
}

type ListGeminiDBInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListGeminiDBBuilder interface {
	ToGeminiDBListDetailQuery() (string, error)
}

func (opts ListGeminiDBInstanceOpts) ToGeminiDBListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListGeminiDBBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToGeminiDBListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageGeminiDBList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GeminiDBPage{pagination.SinglePageBase(r)}
	})

	geminiDBPageheader := map[string]string{"Content-Type": "application/json"}
	pageGeminiDBList.Headers = geminiDBPageheader
	return pageGeminiDBList
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	opts := ListGeminiDBInstanceOpts{
		Id: id,
	}

	page, err := List(client, opts).AllPages()
	r.Body = page.GetBody()
	r.Err = err
	return
}
