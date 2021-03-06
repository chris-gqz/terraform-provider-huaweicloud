package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"
)

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	// TODO: why does back-end convert name to lowercase?
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("huaweicloud_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"huaweicloud_dns_zone_v2.zone_1", "description", "a zone"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_readTTL(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("huaweicloud_dns_zone_v2.zone_1", &zone),
					resource.TestMatchResourceAttr(
						"huaweicloud_dns_zone_v2.zone_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_timeout(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDNS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_timeout(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("huaweicloud_dns_zone_v2.zone_1", &zone),
				),
			},
		},
	})
}

func testAccCheckDNSV2ZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dns_zone_v2" {
			continue
		}

		_, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Zone still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2ZoneExists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
		}

		found, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*zone = *found

		return nil
	}
}

func testAccDNSV2Zone_basic(zoneName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			#type = "PRIMARY"
		}
	`, zoneName)
}

func testAccDNSV2Zone_update(zoneName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			#type = "PRIMARY"
		}
	`, zoneName)
}

func testAccDNSV2Zone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
		}
	`, zoneName)
}

func testAccDNSV2Zone_timeout(zoneName string) string {
	return fmt.Sprintf(`
		resource "huaweicloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email@example.com"
			ttl = 3000

			timeouts {
				create = "5m"
				update = "5m"
				delete = "5m"
			}
		}
	`, zoneName)
}
