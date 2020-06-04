package flavors

import "github.com/huaweicloud/golangsdk"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("flavors")
}
