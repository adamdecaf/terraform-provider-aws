package aws

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSVolumeAttachment_basic(t *testing.T) {
	var i ec2.Instance
	var v ec2.Volume
	resourceName := "aws_volume_attachment.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_name", "/dev/sdh"),
					testAccCheckInstanceExists("aws_instance.test", &i),
					testAccCheckVolumeExists("aws_ebs_volume.test", &v),
					testAccCheckVolumeAttachmentExists(resourceName, &i, &v),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVolumeAttachmentImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSVolumeAttachment_skipDestroy(t *testing.T) {
	var i ec2.Instance
	var v ec2.Volume
	resourceName := "aws_volume_attachment.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfigSkipDestroy(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_name", "/dev/sdh"),
					testAccCheckInstanceExists("aws_instance.test", &i),
					testAccCheckVolumeExists("aws_ebs_volume.test", &v),
					testAccCheckVolumeAttachmentExists(resourceName, &i, &v),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVolumeAttachmentImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"skip_destroy", // attribute only used on resource deletion
				},
			},
		},
	})
}

func TestAccAWSVolumeAttachment_attachStopped(t *testing.T) {
	var i ec2.Instance
	var v ec2.Volume
	resourceName := "aws_volume_attachment.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	stopInstance := func() {
		conn := testAccProvider.Meta().(*AWSClient).ec2conn

		_, err := conn.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: []*string{i.InstanceId},
		})
		if err != nil {
			t.Fatalf("error stopping instance (%s): %s", aws.StringValue(i.InstanceId), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning, ec2.InstanceStateNameStopping},
			Target:     []string{ec2.InstanceStateNameStopped},
			Refresh:    InstanceStateRefreshFunc(conn, *i.InstanceId, []string{}),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			t.Fatalf("Error waiting for instance(%s) to stop: %s", *i.InstanceId, err)
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfigBase(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("aws_instance.test", &i),
				),
			},
			{
				PreConfig: stopInstance,
				Config:    testAccVolumeAttachmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_name", "/dev/sdh"),
					testAccCheckInstanceExists("aws_instance.test", &i),
					testAccCheckVolumeExists("aws_ebs_volume.test", &v),
					testAccCheckVolumeAttachmentExists(resourceName, &i, &v),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVolumeAttachmentImportStateIDFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSVolumeAttachment_update(t *testing.T) {
	resourceName := "aws_volume_attachment.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentUpdateConfig(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "force_detach", "false"),
					resource.TestCheckResourceAttr(resourceName, "skip_destroy", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVolumeAttachmentImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_detach", // attribute only used on resource deletion
					"skip_destroy", // attribute only used on resource deletion
				},
			},
			{
				Config: testAccVolumeAttachmentUpdateConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "force_detach", "true"),
					resource.TestCheckResourceAttr(resourceName, "skip_destroy", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSVolumeAttachmentImportStateIDFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_detach", // attribute only used on resource deletion
					"skip_destroy", // attribute only used on resource deletion
				},
			},
		},
	})
}

func TestAccAWSVolumeAttachment_disappears(t *testing.T) {
	var i ec2.Instance
	var v ec2.Volume
	resourceName := "aws_volume_attachment.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("aws_instance.test", &i),
					testAccCheckVolumeExists("aws_ebs_volume.test", &v),
					testAccCheckVolumeAttachmentExists(resourceName, &i, &v),
					testAccCheckVolumeAttachmentDisappears(resourceName, &i, &v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVolumeAttachmentExists(n string, i *ec2.Instance, v *ec2.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for _, b := range i.BlockDeviceMappings {
			if rs.Primary.Attributes["device_name"] == aws.StringValue(b.DeviceName) {
				if b.Ebs.VolumeId != nil &&
					rs.Primary.Attributes["volume_id"] == aws.StringValue(b.Ebs.VolumeId) &&
					rs.Primary.Attributes["volume_id"] == aws.StringValue(v.VolumeId) {
					// pass
					return nil
				}
			}
		}

		return fmt.Errorf("Error finding instance/volume")
	}
}

func testAccCheckVolumeAttachmentDisappears(n string, i *ec2.Instance, v *ec2.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn

		opts := &ec2.DetachVolumeInput{
			Device:     aws.String(rs.Primary.Attributes["device_name"]),
			InstanceId: i.InstanceId,
			VolumeId:   v.VolumeId,
			Force:      aws.Bool(true),
		}

		_, err := conn.DetachVolume(opts)
		if err != nil {
			return err
		}

		vId := aws.StringValue(v.VolumeId)
		iId := aws.StringValue(i.InstanceId)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{ec2.VolumeAttachmentStateDetaching},
			Target:     []string{ec2.VolumeAttachmentStateDetached},
			Refresh:    volumeAttachmentStateRefreshFunc(conn, rs.Primary.Attributes["device_name"], vId, iId),
			Timeout:    5 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Detaching Volume (%s) from Instance (%s)", vId, iId)
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for Volume (%s) to detach from Instance (%s): %s",
				vId, iId, err)
		}

		return err
	}
}

func testAccCheckVolumeAttachmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ec2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_volume_attachment" {
			continue
		}

		request := &ec2.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(rs.Primary.Attributes["volume_id"])},
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("attachment.device"),
					Values: []*string{aws.String(rs.Primary.Attributes["device_name"])},
				},
				{
					Name:   aws.String("attachment.instance-id"),
					Values: []*string{aws.String(rs.Primary.Attributes["instance_id"])},
				},
			},
		}

		_, err := conn.DescribeVolumes(request)
		if err != nil {
			if isAWSErr(err, "InvalidVolume.NotFound", "") {
				return nil
			}
			return fmt.Errorf("error describing volumes (%s): %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccVolumeAttachmentInstanceOnlyConfigBase(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_ami" "amzn-ami-minimal-hvm-ebs" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-minimal-hvm-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

data "aws_ec2_instance_type_offering" "available" {
  filter {
    name   = "instance-type"
    values = ["t3.micro", "t2.micro"]
  }

  location_type            = "availability-zone"
  preferred_instance_types = ["t3.micro", "t2.micro"]
}

resource "aws_instance" "test" {
  ami               = "${data.aws_ami.amzn-ami-minimal-hvm-ebs.id}"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  instance_type     = "${data.aws_ec2_instance_type_offering.available.instance_type}"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVolumeAttachmentConfigBase(rName string) string {
	return testAccVolumeAttachmentInstanceOnlyConfigBase(rName) + fmt.Sprintf(`
resource "aws_ebs_volume" "test" {
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  size              = 1

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVolumeAttachmentConfig(rName string) string {
	return testAccVolumeAttachmentConfigBase(rName) + fmt.Sprintf(`
resource "aws_volume_attachment" "test" {
  device_name = "/dev/sdh"
  volume_id   = "${aws_ebs_volume.test.id}"
  instance_id = "${aws_instance.test.id}"
}
`)
}

func testAccVolumeAttachmentConfigSkipDestroy(rName string) string {
	return testAccVolumeAttachmentConfigBase(rName) + fmt.Sprintf(`
data "aws_ebs_volume" "test" {
  filter {
    name   = "size"
    values = ["${aws_ebs_volume.test.size}"]
  }
  filter {
    name   = "availability-zone"
    values = ["${aws_ebs_volume.test.availability_zone}"]
  }
  filter {
    name   = "tag:Name"
    values = ["%[1]s"]
  }
}

resource "aws_volume_attachment" "test" {
  device_name  = "/dev/sdh"
  volume_id    = "${data.aws_ebs_volume.test.id}"
  instance_id  = "${aws_instance.test.id}"
  skip_destroy = true
}
`, rName)
}

func testAccVolumeAttachmentUpdateConfig(rName string, detach bool) string {
	return testAccVolumeAttachmentConfigBase(rName) + fmt.Sprintf(`
resource "aws_volume_attachment" "test" {
  device_name  = "/dev/sdh"
  volume_id    = "${aws_ebs_volume.test.id}"
  instance_id  = "${aws_instance.test.id}"
  force_detach = %[1]t
  skip_destroy = %[1]t
}
`, detach)
}

func testAccAWSVolumeAttachmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["device_name"], rs.Primary.Attributes["volume_id"], rs.Primary.Attributes["instance_id"]), nil
	}
}
