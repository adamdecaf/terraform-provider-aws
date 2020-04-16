package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsEip_Filter(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigFilter(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_dns", resourceName, "public_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_ip", resourceName, "public_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_Id(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_dns", resourceName, "public_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_ip", resourceName, "public_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_PublicIP_EC2Classic(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	// Do not parallelize this test until the provider testing framework
	// has a stable us-east-1 alias
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigPublicIpEc2Classic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_dns", resourceName, "public_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_ip", resourceName, "public_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_PublicIP_VPC(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigPublicIpVpc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_dns", resourceName, "public_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_ip", resourceName, "public_ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domain", resourceName, "domain"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_Tags(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigTags(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_dns", resourceName, "public_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "public_ip", resourceName, "public_ip"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_NetworkInterface(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "network_interface_id", resourceName, "network_interface"),
					resource.TestCheckResourceAttrPair(dataSourceName, "private_dns", resourceName, "private_dns"),
					resource.TestCheckResourceAttrPair(dataSourceName, "private_ip", resourceName, "private_ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domain", resourceName, "domain"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsEip_Instance(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsEipConfigInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instance_id", resourceName, "instance"),
					resource.TestCheckResourceAttrPair(dataSourceName, "association_id", resourceName, "association_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAWSEIP_COIPPool(t *testing.T) {
	dataSourceName := "data.aws_eip.test"
	resourceName := "aws_eip.test"

	poolId := os.Getenv("AWS_COIP_POOL_ID")
	if poolId == "" {
		t.Skip(
			"Environment variable AWS_COIP_POOL_ID is not set. " +
				"This environment variable must be set to the ID of " +
				"a deployed Coip Pool to enable this test.")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSEIPConfig_CoipPool(poolId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "customer_owned_ipv4_pool", dataSourceName, "customer_owned_ipv4_pool"),
					resource.TestCheckResourceAttrSet(dataSourceName, "customer_owned_ip"),
				),
			},
		},
	})
}

func testAccDataSourceAWSEIPConfig_CoipPool(poolId string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  vpc                      = true
  customer_owned_ipv4_pool = "%s"
}

data "aws_eip" "test" {
  id = "${aws_eip.test.id}"
}
`, poolId)
}

func testAccDataSourceAwsEipConfigFilter(rName string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  vpc = true

  tags = {
    Name = %q
  }
}

data "aws_eip" "test" {
  filter {
    name   = "tag:Name"
    values = ["${aws_eip.test.tags.Name}"]
  }
}
`, rName)
}

const testAccDataSourceAwsEipConfigId = `
resource "aws_eip" "test" {
  vpc = true
}

data "aws_eip" "test" {
  id = "${aws_eip.test.id}"
}
`

const testAccDataSourceAwsEipConfigPublicIpEc2Classic = `
provider "aws" {
  region = "us-east-1"
}

resource "aws_eip" "test" {}

data "aws_eip" "test" {
  public_ip = "${aws_eip.test.public_ip}"
}
`

const testAccDataSourceAwsEipConfigPublicIpVpc = `
resource "aws_eip" "test" {
  vpc = true
}

data "aws_eip" "test" {
  public_ip = "${aws_eip.test.public_ip}"
}
`

func testAccDataSourceAwsEipConfigTags(rName string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  vpc = true

  tags = {
    Name = %q
  }
}

data "aws_eip" "test" {
  tags = {
    Name = "${aws_eip.test.tags["Name"]}"
  }
}
`, rName)
}

const testAccDataSourceAwsEipConfigNetworkInterface = `
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "test" {
  vpc_id = "${aws_vpc.test.id}"
  cidr_block = "10.1.0.0/24"
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_network_interface" "test" {
  subnet_id = "${aws_subnet.test.id}"
}

resource "aws_eip" "test" {
  vpc = true
  network_interface = "${aws_network_interface.test.id}"
}

data "aws_eip" "test" {
  filter {
    name   = "network-interface-id"
    values = ["${aws_eip.test.network_interface}"]
  }
}
`

const testAccDataSourceAwsEipConfigInstance = `
data "aws_availability_zones" "available" {
  # Error launching source instance: Unsupported: Your requested instance type (t2.micro) is not supported in your requested Availability Zone (us-west-2d).
  blacklisted_zone_ids = ["usw2-az4"]
  state                = "available"
}

resource "aws_vpc" "test" {
  cidr_block = "10.2.0.0/16"
}

resource "aws_subnet" "test" {
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  vpc_id = "${aws_vpc.test.id}"
  cidr_block = "10.2.0.0/24"
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"
}

data "aws_ami" "test" {
  most_recent = true
  name_regex  = "^amzn-ami.*ecs-optimized$"

  owners = [
    "amazon",
  ]
}

resource "aws_instance" "test" {
  ami = "${data.aws_ami.test.id}"
  subnet_id = "${aws_subnet.test.id}"
  instance_type = "t2.micro"
}

resource "aws_eip" "test" {
  vpc = true
  instance = "${aws_instance.test.id}"
}

data "aws_eip" "test" {
  filter {
    name = "instance-id"
    values = ["${aws_eip.test.instance}"]
  }
}
`
