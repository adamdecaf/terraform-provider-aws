package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsSesSendingIpPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSesSendingIpPoolCreate,
		Read:   resourceAwsSesSendingIpPoolRead,
		Delete: resourceAwsSesSendingIpPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAwsSesSendingIpPoolCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesv2Conn

	poolName := d.Get("name").(string)
	_, err := conn.CreateDedicatedIpPool(&sesv2.CreateDedicatedIpPoolInput{
		PoolName: aws.String(poolName),
	})
	if err != nil {
		return fmt.Errorf("Error creating SESv2 ip pool: %s", err)
	}

	d.SetId(poolName)

	return resourceAwsSesSendingIpPoolRead(d, meta)
}

func resourceAwsSesSendingIpPoolRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesv2Conn

	response, err := conn.ListDedicatedIpPools(&sesv2.ListDedicatedIpPoolsInput{
		PageSize: aws.Int64(100),
	})
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error reading SESv2 ip pool: %v", err)
	}
	for i := range response.DedicatedIpPools {
		if n := *response.DedicatedIpPools[i]; strings.EqualFold(n, d.Id()) {
			return nil
		}
	}
	return nil
}

func resourceAwsSesSendingIpPoolDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesv2Conn
	log.Printf("[DEBUG] SES Delete Sending IP Pool: id=%s name=%s", d.Id(), d.Get("name").(string))
	_, err := conn.DeleteDedicatedIpPool(&sesv2.DeleteDedicatedIpPoolInput{
		PoolName: aws.String(d.Id()),
	})
	return err
}
