package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTeamPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamPipelineCreate,
		Read:   resourceTeamPipelineRead,
		Update: resourceTeamPipelineUpdate,
		Delete: resourceTeamPipelineDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTeamPipelineCreate(d *schema.ResourceData, m interface{}) error {
	return resourceTeamPipelineRead(d, m)
}

func resourceTeamPipelineRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceTeamPipelineUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceTeamPipelineRead(d, m)
}

func resourceTeamPipelineDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
