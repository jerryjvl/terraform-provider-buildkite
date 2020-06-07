package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamMemberCreate,
		Read:   resourceTeamMemberRead,
		Update: resourceTeamMemberUpdate,
		Delete: resourceTeamMemberDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTeamMemberCreate(d *schema.ResourceData, m interface{}) error {
	return resourceTeamMemberRead(d, m)
}

func resourceTeamMemberRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTeamMemberUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTeamMemberRead(d, m)
}

func resourceTeamMemberDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
