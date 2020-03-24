package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAwsSESConfigurationSetDeliveryOptions_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSSES(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsSESConfigurationSetDeliveryOptionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsSESConfigurationSetDeliveryOptionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsSESConfigurationSetDeliveryOptionsExists("aws_ses_configuration_set_deliver_options.test"),
				),
			},
		},
	})
}

func testAccCheckAwsSESConfigurationSetDeliveryOptionsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).sesv2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ses_configuration_set_deliver_options" {
			continue
		}

		response, err := conn.ListDedicatedIpPools(&sesv2.ListDedicatedIpPoolsInput{})
		if err != nil {
			return err
		}

		found := false
		for i := range response.DedicatedIpPools {
			if n := *response.DedicatedIpPools[i]; n == "sender" {
				found = true
			}
		}
		if found {
			return errors.New("The configuration set delivery options still exists")
		}
	}

	return nil
}

func testAccCheckAwsSESConfigurationSetDeliveryOptionsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).sesv2Conn

		response, err := conn.ListConfigurationSets(&sesv2.ListConfigurationSetsInput{})
		if err != nil {
			return err
		}

		found := false
		for i := range response.ConfigurationSets {
			if *response.ConfigurationSets[i] == "test" {
				found = true
			}
		}

		if !found {
			return errors.New("The configuration set delivery options was not created")
		}

		return nil
	}
}

const testAccAwsSESConfigurationSetDeliveryOptionsConfig = `
resource "aws_ses_sending_ip_pool" "test" {
  name = "sender"
}
resource "aws_ses_configuration_set" "test" {
  name = "test"
}
resource "aws_ses_configuration_set_delivery_options" "test" {
  configuration_set = "test"
  sending_pool      = "sender"

  depends_on = ["aws_ses_sending_ip_pool.test"]
}
`
