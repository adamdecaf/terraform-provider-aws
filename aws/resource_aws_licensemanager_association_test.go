package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/licensemanager"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSLicenseManagerAssociation_basic(t *testing.T) {
	var licenseSpecification licensemanager.LicenseSpecification
	resourceName := "aws_licensemanager_association.example"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLicenseManagerAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseManagerAssociationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLicenseManagerAssociationExists(resourceName, &licenseSpecification),
					resource.TestCheckResourceAttrPair(resourceName, "license_configuration_arn", "aws_licensemanager_license_configuration.example", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "resource_arn", "aws_instance.example", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLicenseManagerAssociationExists(resourceName string, licenseSpecification *licensemanager.LicenseSpecification) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		resourceArn, licenseConfigurationArn, err := parseLicenseManagerAssociationId(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*AWSClient).licensemanagerconn
		resp, err := conn.ListLicenseSpecificationsForResource(&licensemanager.ListLicenseSpecificationsForResourceInput{
			ResourceArn: aws.String(resourceArn),
		})

		if err != nil {
			return fmt.Errorf("Error retrieving License Manager association (%s): %s", rs.Primary.ID, err)
		}

		for _, ls := range resp.LicenseSpecifications {
			if aws.StringValue(ls.LicenseConfigurationArn) == licenseConfigurationArn {
				*licenseSpecification = *ls
				return nil
			}
		}

		return fmt.Errorf("Error retrieving License Manager association (%s): Not found", rs.Primary.ID)
	}
}

func testAccCheckLicenseManagerAssociationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).licensemanagerconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_licensemanager_association" {
			continue
		}

		// Try to find the resource
		resourceArn, licenseConfigurationArn, err := parseLicenseManagerAssociationId(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.ListLicenseSpecificationsForResource(&licensemanager.ListLicenseSpecificationsForResourceInput{
			ResourceArn: aws.String(resourceArn),
		})

		if err != nil {
			return err
		}

		for _, ls := range resp.LicenseSpecifications {
			if aws.StringValue(ls.LicenseConfigurationArn) == licenseConfigurationArn {
				return fmt.Errorf("License Manager association %q still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccLicenseManagerAssociationConfig_basic = `
data "aws_ami" "example" {
  most_recent      = true

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }

  filter {
    name   = "name"
    values = ["amzn-ami-vpc-nat*"]
  }
}

resource "aws_instance" "example" {
  ami           = "${data.aws_ami.example.id}"
  instance_type = "t2.micro"
}

resource "aws_licensemanager_license_configuration" "example" {
  name                  = "Example"
  license_counting_type = "vCPU"
}

resource "aws_licensemanager_association" "example" {
  license_configuration_arn = "${aws_licensemanager_license_configuration.example.id}"
  resource_arn              = "${aws_instance.example.arn}"
}
`
