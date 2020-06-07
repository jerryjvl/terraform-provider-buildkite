package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSsoProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceSsoProviderCreate,
		Read:   resourceSsoProviderRead,
		Update: resourceSsoProviderUpdate,
		Delete: resourceSsoProviderDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSsoProviderCreate(d *schema.ResourceData, m interface{}) error {
	return resourceSsoProviderRead(d, m)
}

func resourceSsoProviderRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSsoProviderUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSsoProviderRead(d, m)
}

func resourceSsoProviderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
