package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for builds.
func dataSourceBuilds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBuildsRead,

		Schema: map[string]*schema.Schema{
			"builds": schemaBuildList(DataSourceFullEntity),
		},
	}
}

// Read builds from the Buildkite API and convert to the Terraform schema.
func dataSourceBuildsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	builds, err := client.readBuilds()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range builds {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("builds", list); err != nil {
		return fmt.Errorf("error setting builds: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of builds.
func schemaBuildList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite builds",
		Elem:        schemaBuild(mode),
	}
}

// Construct a Terraform schema definition for a build.
func schemaBuild(mode SchemaMode) *schema.Resource {
	switch mode {
	case DataSourceReferenceOnly:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"uuid": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	case DataSourceFullEntity:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"uuid": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	default:
		return &schema.Resource{}
	}
}

// Build defines the properties on the Buildkite API to map to Terraform.
type Build struct {
	Annotations   []Annotation
	Branch        string
	CanceledAt    string
	CanceledBy    User
	Commit        string
	CreatedAt     string
	CreatedBy     User
	Env           []string
	FinishedAt    string
	ID            string
	Message       string
	Number        int
	Pipeline      Pipeline
	PullRequest   struct{ ID string }
	RebuiltFrom   string
	ScheduledAt   string
	Source        struct{ Name string }
	StartedAt     string
	State         string
	TriggeredFrom struct {
		Build string
		ID    string
		UUID  string
	}
	URL  string
	UUID string
}

type Annotation struct {
	Body struct {
		HTML string
		Text string
	}
	Context   string
	CreatedAt string
	ID        string
	Style     string
	UpdatedAt string
	UUID      string
}

/*
type Artifact struct {
	DownloadURL graphql.String
	ID          graphql.ID
	MimeType    graphql.String
	Path        graphql.String
	Sha1Sum     graphql.String
	Size        graphql.Int
	State       graphql.String
	UUID        graphql.ID
}
*/

// BuildList defines the properties on the Buildkite API to map to Terraform.
type BuildList struct {
	Count int
	Edges []struct {
		Node Build
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Build) convert(mode SchemaMode) map[string]interface{} {
	switch mode {
	case DataSourceReferenceOnly:
		return map[string]interface{}{
			"uuid": source.UUID,
		}
	case DataSourceFullEntity:
		return map[string]interface{}{
			"uuid": source.UUID,
		}
	default:
		return map[string]interface{}{}
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *BuildList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readBuilds() ([]Build, error) {
	return nil, nil
}
