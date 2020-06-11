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
					testAccCheckGeminiDBV3InstanceExists("huaweicloud_geminidb_instance.instance_5", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_geminidb_instance.instance_5", "name", "geminidb_instance_5"),
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
provider "huaweicloud" {
    user_name   = "terraform"
    password    = "OpenStackSDK@123"
    domain_name = "freesky-edward"
    #tenant_id   = "e71eaa2d4efc4567abf35458e0f504da"
    tenant_name = "cn-south-1"
    region      = "cn-south-1"
    auth_url    = "https://iam.myhwclouds.com:443/v3"
}

resource "huaweicloud_vpc_v1" "vpc" {
  name = "terraform_vpc_v1_test_5"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_networking_network_v2" "network_5" {
  name = "network_5"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_5" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = huaweicloud_networking_network_v2.network_5.id
}

resource "huaweicloud_networking_secgroup_v2" "secgroup_5" {
  name = "secgroup_5"
}

resource "huaweicloud_geminidb_instance" "instance_5" {
  name        = "geminidb_instance_5"
  flavor      = "geminidb.cassandra.xlarge.4"
  password    = "OpenStackSDK@123"
  volume_size = 100
  vpc_id      = huaweicloud_vpc_v1.vpc.id
  subnet_id   = huaweicloud_networking_subnet_v2.subnet_5.id
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup_5.id
  availability_zone = "cn-south-2b"
}
`)

func testAccCheckGeminiDBV3InstanceExists(n string, instance *instances.CreateGeminiDB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.GeminiDBV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s ", err)
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
		if instances.TotalCount == 0 {
			return fmt.Errorf("Instance not found. ")
		}
		return nil
	}
}
