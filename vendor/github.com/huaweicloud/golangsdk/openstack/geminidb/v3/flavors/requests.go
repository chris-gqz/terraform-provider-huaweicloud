package flavors

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type AllInstancesFlavor struct {
	Region     string `q:"region" required:"true"`
	EngineName string `q:"engine_name"`
}

type AllInstancesFlavorBuilder interface {
	ToAllInstancesFlavorDetailQuery() (string, error)
}

func (opts AllInstancesFlavor) ToAllInstancesFlavorDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func GetAllInstancesFlavor(client *golangsdk.ServiceClient, opts AllInstancesFlavorBuilder) pagination.Pager {
	url := getURL(client)
	if opts != nil {
		query, err := opts.ToAllInstancesFlavorDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageAllInstancesFlavor := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AllInstancesFlavorPage{pagination.SinglePageBase(r)}
	})

	allInstancesFlavorPageheader := map[string]string{"Content-Type": "application/json"}
	pageAllInstancesFlavor.Headers = allInstancesFlavorPageheader
	return pageAllInstancesFlavor
}
