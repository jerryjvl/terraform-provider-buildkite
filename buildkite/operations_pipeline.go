package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for pipelines.
func dataSourcePipelines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePipelinesRead,

		Schema: map[string]*schema.Schema{
			"pipelines": schemaPipelineList(DataSourceFullEntity),
		},
	}
}

// Read pipelines from the Buildkite API and convert to the Terraform schema.
func dataSourcePipelinesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	pipelines, err := client.readPipelines()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range pipelines {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("pipelines", list); err != nil {
		return fmt.Errorf("error setting pipelines: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of pipelines.
func schemaPipelineList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite pipelines",
		Elem:        schemaPipeline(mode),
	}
}

// Construct a Terraform schema definition for a pipeline.
func schemaPipeline(mode SchemaMode) *schema.Resource {
	switch mode {
	case DataSourceReferenceOnly:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"slug": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	case DataSourceFullEntity:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"slug": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	default:
		return &schema.Resource{}
	}
}

// Pipeline defines the properties on the Buildkite API to map to Terraform.
type Pipeline struct {
	CancelIntermediateBuilds             bool
	CancelIntermediateBuildsBranchFilter string
	CommitShortLength                    int
	CreatedAt                            string
	DefaultBranch                        string
	Description                          string
	Favorite                             bool
	ID                                   string
	Name                                 string
	NextBuildNumber                      int
	Repository                           struct {
		Provider struct {
			Name       string
			URL        string
			WebhookURL string
		}
		URL string
	}
	Schedules                          []PipelineSchedule
	SkipIntermediateBuilds             bool
	SkipIntermediateBuildsBranchFilter string
	Slug                               string
	Steps                              struct{ YAML string }
	URL                                string
	UUID                               string
	Visibility                         string
	WebhookURL                         string
}

type PipelineSchedule struct {
	Branch        string
	Commit        string
	CreatedAt     string
	CreatedBy     User
	Cronline      string
	Enabled       bool
	Env           []string
	FailedAt      string
	FailedMessage string
	ID            string
	Label         string
	Message       string
	NextBuildAt   string
	OwnedBy       User
	Pipeline      Pipeline
	UUID          string
}

// PipelineList defines the properties on the Buildkite API to map to Terraform.
type PipelineList struct {
	Count int
	Edges []struct {
		Node Pipeline
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Pipeline) convert(mode SchemaMode) map[string]interface{} {
	switch mode {
	case DataSourceReferenceOnly:
		return map[string]interface{}{
			"slug": source.Slug,
		}
	case DataSourceFullEntity:
		return map[string]interface{}{
			"slug": source.Slug,
		}
	default:
		return map[string]interface{}{}
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *PipelineList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readPipelines() ([]Pipeline, error) {
	return nil, nil
}
