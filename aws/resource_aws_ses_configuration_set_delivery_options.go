package aws

import (
	"fmt"
	"log"
	"strings"

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
	configSet, err := findSesV2ConfigurationSet(d, d.Id(), nil, meta)

	if configSet == nil {
		log.Printf("[WARN] SES Configuration Set (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.Set("name", d.Id())
	d.Set("configuration_set", d.Id())

	if opts := configSet.DeliveryOptions; opts != nil {
		d.Set("tls_policy", opts.TlsPolicy)
		d.Set("sending_pool", opts.SendingPoolName)
	}

	return nil
}

func resourceAwsSesConfigurationSetDeliveryOptionsUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAwsSesConfigurationSetDeliveryOptionsCreate(d, meta)
}

func resourceAwsSesConfigurationSetDeliveryOptionsDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func findSesV2ConfigurationSet(d *schema.ResourceData, name string, token *string, meta interface{}) (*sesv2.GetConfigurationSetOutput, error) {
	conn := meta.(*AWSClient).sesv2Conn

	listOpts := &sesv2.ListConfigurationSetsInput{
		NextToken: token,
	}

	configurationSetExists := false
	response, err := conn.ListConfigurationSets(listOpts)
	for i := range response.ConfigurationSets {
		if strings.EqualFold(*response.ConfigurationSets[i], name) {
			configurationSetExists = true
		}
	}
	if err != nil && !configurationSetExists && response.NextToken != nil {
		return findSesV2ConfigurationSet(d, name, response.NextToken, meta)
	}
	if err != nil {
		return nil, err
	}

	configSet, err := conn.GetConfigurationSet(&sesv2.GetConfigurationSetInput{
		ConfigurationSetName: aws.String(name),
	})

	return configSet, err
}
