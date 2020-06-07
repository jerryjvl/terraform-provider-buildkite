package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for teams.
func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTeamsRead,

		Schema: map[string]*schema.Schema{
			"teams": schemaTeamList(DataSourceFullEntity),
		},
	}
}

// Read teams from the Buildkite API and convert to the Terraform schema.
func dataSourceTeamsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	teams, err := client.readTeams()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range teams {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("teams", list); err != nil {
		return fmt.Errorf("error setting teams: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of teams.
func schemaTeamList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite teams",
		Elem:        schemaTeam(mode),
	}
}

// Construct a Terraform schema definition for a team.
func schemaTeam(mode SchemaMode) *schema.Resource {
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

// Team defines the properties on the Buildkite API to map to Terraform.
type Team struct {
	CreatedAt         string
	CreatedBy         User
	DefaultMemberRole string
	Description       string
	ID                string
	IsDefaultTeam     bool
	//Members                   []TeamMember
	MembersCanCreatePipelines bool
	Name                      string
	//Pipelines                 []TeamPipeline
	Privacy string
	Slug    string
	UUID    string
}

/*
type TeamMember struct {
	CreatedAt graphql.String
	CreatedBy User
	ID        graphql.ID
	Role      graphql.String
	Team      Team
	User      User
	UUID      graphql.ID
}

type TeamPipeline struct {
	AccessLevel graphql.String
	CreatedAt   graphql.String
	CreatedBy   User
	ID          graphql.ID
	Pipeline    Pipeline
	Team        Team
	UUID        graphql.ID
}
*/

// TeamList defines the properties on the Buildkite API to map to Terraform.
type TeamList struct {
	Count int
	Edges []struct {
		Node Team
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Team) convert(mode SchemaMode) map[string]interface{} {
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
func (source *TeamList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readTeams() ([]Team, error) {
	return nil, nil
}
