package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
)

func TestGeminiDBInstance_basic(t *testing.T) {
	var instance instances.CreateGeminiDB

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGeminiDBV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccGeminiDBInstanceV3Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists("huaweicloud_dds_instance_v3.instance", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_dds_instance_v3.instance", "name", "dds-instance"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.GeminiDBV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_geminidb_instance" {
			continue
		}

		opts := instances.ListGeminiDBInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractGeminiDBInstances(allPages)
		if err != nil {
			return err
		}

		if instances.TotalCount > 0 {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

var TestAccGeminiDBInstanceV3Config_basic = fmt.Sprintf(`
resource "huaweicloud_geminidb_instance" "instance_1" {
 name        = "geminidb_instance_1"
 flavor      = "geminidb.cassandra.xlarge.4"
 password    = var.password
 volume_size = 100
 vpc_id      = "%s"
 subnet_id   = "%s"
 security_group_id = var.secgroup_id
 availability_zone = "%s"

 backup_strategy {
   start_time = "03:00-05:00"
   keep_days  = 14
 }
}
`, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
