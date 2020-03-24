package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsSesConfigurationSetDeliveryOptions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSesConfigurationSetDeliveryOptionsCreate,
		Read:   resourceAwsSesConfigurationSetDeliveryOptionsRead,
		Update: resourceAwsSesConfigurationSetDeliveryOptionsUpdate,
		Delete: resourceAwsSesConfigurationSetDeliveryOptionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"configuration_set": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sending_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsSesConfigurationSetDeliveryOptionsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesv2Conn

	name := d.Get("configuration_set").(string)
	opts := &sesv2.PutConfigurationSetDeliveryOptionsInput{
		ConfigurationSetName: aws.String(name),
	}
	if v, ok := d.GetOk("sending_pool"); ok {
		opts.SendingPoolName = aws.String(v.(string))
	}
	if v, ok := d.GetOk("tls_policy"); ok {
		opts.TlsPolicy = aws.String(v.(string))
	}

	_, err := conn.PutConfigurationSetDeliveryOptions(opts)
	if err != nil {
		return fmt.Errorf("Error creating SES configuration set delivery options: %s", err)
	}

	d.SetId(name)

	return resourceAwsSesConfigurationSetDeliveryOptionsRead(d, meta)
}

func resourceAwsSesConfigurationSetDeliveryOptionsRead(d *schema.ResourceData, meta interface{}) error {
	configurationSetExists, err := findSesV2ConfigurationSet(d.Id(), nil, meta)

	if !configurationSetExists {
		log.Printf("[WARN] SES Configuration Set (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", d.Id())

	return nil
}

func resourceAwsSesConfigurationSetDeliveryOptionsUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAwsSesConfigurationSetDeliveryOptionsCreate(d, meta)
}

func resourceAwsSesConfigurationSetDeliveryOptionsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func findSesV2ConfigurationSet(name string, token *string, meta interface{}) (bool, error) {
	conn := meta.(*AWSClient).sesv2Conn

	configurationSetExists := false

	listOpts := &sesv2.ListConfigurationSetsInput{
		NextToken: token,
	}

	response, err := conn.ListConfigurationSets(listOpts)
	for i := range response.ConfigurationSets {
		if *response.ConfigurationSets[i] == name {
			configurationSetExists = true
		}
	}

	if err != nil && !configurationSetExists && response.NextToken != nil {
		configurationSetExists, err = findSesV2ConfigurationSet(name, response.NextToken, meta)
	}

	if err != nil {
		return false, err
	}

	return configurationSetExists, nil
}
