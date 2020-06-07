package buildkite

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for agents.
func dataSourceAgents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentsRead,

		Schema: map[string]*schema.Schema{
			"organization_slug": {
				Type:        schema.TypeString,
				Description: "Buildkite organization slug",
				Optional:    true,
			},

			"agents": schemaAgentList(DataSourceFullEntity),
		},
	}
}

// Read agents from the Buildkite API and convert to the Terraform schema.
func dataSourceAgentsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	slug := d.Get("organization_slug").(string)

	agents, err := client.readAgents(slug)
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range agents {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("agents", list); err != nil {
		return fmt.Errorf("error setting agents: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of agents.
func schemaAgentList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite agents",
		Elem:        schemaAgent(mode),
	}
}

// Construct a Terraform schema definition for an agent.
func schemaAgent(mode SchemaMode) *schema.Resource {
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
				"hostname": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"id": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"ip_address": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"is_deprecated": &schema.Schema{
					Type:     schema.TypeBool,
					Computed: true,
				},
				"meta_data": &schema.Schema{
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"operating_system": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"public": &schema.Schema{
					Type:     schema.TypeBool,
					Computed: true,
				},
				"user_agent": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"uuid": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"version": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"version_has_known_issues": &schema.Schema{
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		}
	default:
		return &schema.Resource{}
	}
}

// Agent defines the properties on the Buildkite API to map to Terraform.
type Agent struct {
	Hostname              string
	ID                    string
	IPAddress             string
	IsDeprecated          bool
	MetaData              []string
	Name                  string
	OperatingSystem       struct{ Name string }
	Public                bool
	UserAgent             string
	UUID                  string
	Version               string
	VersionHasKnownIssues bool
}

// AgentList defines the properties on the Buildkite API to map to Terraform.
type AgentList struct {
	Count int
	Edges []struct {
		Node Agent
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Agent) convert(mode SchemaMode) map[string]interface{} {
	switch mode {
	case DataSourceReferenceOnly:
		return map[string]interface{}{
			"uuid": source.UUID,
		}
	case DataSourceFullEntity:
		metaData := make(map[string]string)
		for _, e := range source.MetaData {
			parts := strings.Split(string(e), "=")
			metaData[parts[0]] = parts[1]
		}

		return map[string]interface{}{
			"hostname":                 source.Hostname,
			"id":                       source.ID,
			"ip_address":               source.IPAddress,
			"is_deprecated":            source.IsDeprecated,
			"meta_data":                metaData,
			"name":                     source.Name,
			"operating_system":         source.OperatingSystem.Name,
			"public":                   source.Public,
			"user_agent":               source.UserAgent,
			"uuid":                     source.UUID,
			"version":                  source.Version,
			"version_has_known_issues": source.VersionHasKnownIssues,
		}
	default:
		return map[string]interface{}{}
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *AgentList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readAgents(slug string) ([]Agent, error) {
	return nil, nil
	/*
		coalescedSlug := slug
		if coalescedSlug == "" {
			coalescedSlug = client.slug
		}

		var countQuery struct {
			Organization struct {
				Agents struct {
					Count graphql.Int
				}
			} `graphql:"organization(slug: $slug)"`
		}
		countVars := map[string]interface{}{
			"slug": graphql.ID(coalescedSlug),
		}
		err := client.graph.Query(context.Background(), &countQuery, countVars)
		if err != nil {
			return nil, err
		}

		var query struct {
			Organization struct {
				Agents struct {
					Edges []struct {
						Node Agent
					}
				} `graphql:"agents(first: $first)"`
			} `graphql:"organization(slug: $slug)"`
		}
		vars := map[string]interface{}{
			"slug":  graphql.ID(coalescedSlug),
			"first": countQuery.Organization.Agents.Count,
		}

		err = client.graphQL.Query(context.Background(), &query, vars)
		if err != nil {
			return nil, err
		}

		var agents []Agent
		edges := &query.Organization.Agents.Edges
		for _, agent := range *edges {
			agents = append(agents, agent.Node)
		}
		return agents, nil
	*/
}
