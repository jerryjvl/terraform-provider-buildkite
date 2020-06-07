package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePipelineSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineScheduleCreate,
		Read:   resourcePipelineScheduleRead,
		Update: resourcePipelineScheduleUpdate,
		Delete: resourcePipelineScheduleDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePipelineScheduleCreate(d *schema.ResourceData, m interface{}) error {
	return resourcePipelineScheduleRead(d, m)
}

func resourcePipelineScheduleRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	return resourcePipelineScheduleRead(d, m)
}

func resourcePipelineScheduleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
