package buildkite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/peterhellberg/link"
)

// Define a Terraform data source for organizations.
func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationsRead,

		Schema: map[string]*schema.Schema{
			"organizations": schemaOrganizationList(DataSourceFullEntity),
		},
	}
}

// Read organizations from the Buildkite API and convert to the Terraform schema.
func dataSourceOrganizationsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	organizations, err := client.readOrganizations()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range organizations {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("organizations", list); err != nil {
		return fmt.Errorf("error setting organizations: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of organizations.
func schemaOrganizationList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite organizations",
		Elem:        schemaOrganization(mode),
	}
}

// Construct a Terraform schema definition for an organization.
func schemaOrganization(mode SchemaMode) *schema.Resource {
	switch mode {
	case DataSourceReferenceOnly:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"slug": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	case DataSourceFullEntity:
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agents": schemaAgentList(DataSourceReferenceOnly),
				"icon_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"members": schemaMemberList(DataSourceReferenceOnly),
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"pipelines": schemaPipelineList(DataSourceReferenceOnly),
				"public": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"slug": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sso_providers": schemaSsoProviderList(DataSourceReferenceOnly),
				"teams":         schemaTeamList(DataSourceReferenceOnly),
				"url": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"uuid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"web_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		}
	default:
		return &schema.Resource{}
	}
}

// Organization defines the properties on the Buildkite API to map to Terraform.
type Organization struct {
	Agents       AgentList
	IconURL      string
	ID           string
	Members      MemberList
	Name         string
	Pipelines    PipelineList
	Public       bool
	Slug         string
	SsoProviders SsoProviderList
	Teams        TeamList
	URL          string
	UUID         string
	WebURL       string `json:"web_url"`
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Organization) convert(mode SchemaMode) map[string]interface{} {
	switch mode {
	case DataSourceReferenceOnly:
		return map[string]interface{}{
			"slug": source.Slug,
		}
	case DataSourceFullEntity:
		return map[string]interface{}{
			"agents":        source.Agents.convert(DataSourceReferenceOnly),
			"icon_url":      source.IconURL,
			"id":            source.ID,
			"members":       source.Members.convert(DataSourceReferenceOnly),
			"name":          source.Name,
			"pipelines":     source.Pipelines.convert(DataSourceReferenceOnly),
			"public":        source.Public,
			"slug":          source.Slug,
			"sso_providers": source.SsoProviders.convert(DataSourceReferenceOnly),
			"teams":         source.Teams.convert(DataSourceReferenceOnly),
			"url":           source.URL,
			"uuid":          source.UUID,
			"web_url":       source.WebURL,
		}
	default:
		return map[string]interface{}{}
	}
}

// Retrieve the basic details about the organization and lengths of sub-lists to query.
const queryOrganizationCount = "organization(slug: %s) { " +
	"agents { count } " +
	"iconUrl " +
	"id " +
	"members { count } " +
	"name pipelines { count } " +
	"public " +
	"slug " +
	"ssoProviders { count } " +
	"teams { count } " +
	"uuid }"

// Retrieve the key identifiers for sub-lists.
const queryOrganizationLists = "organization(slug: %s) { " +
	"agents(first: %d) { edges { node { uuid } } } " +
	"members(first: %d) { edges { node { uuid } } } " +
	"pipelines(first: %d) { edges { node { slug } } } " +
	"ssoProviders(first: %d) { edges { node { uuid } } } " +
	"teams(first: %d) { edges { node { uuid } } } }"

// Read a list of all defined organizations from the Buildkite API.
func (client *Client) readOrganizations() ([]Organization, error) {
	// Get Buildkite organizations until there are no next pages
	var organizations []Organization
	url := "https://api.buildkite.com/v2/organizations"
	for {
		res, err := client.httpAPI.Get(url)
		if err != nil {
			return nil, err
		}

		var temp []Organization
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bodyBytes, &temp); err != nil {
			return nil, err
		}

		organizations = append(organizations, temp...)

		link, linkOk := link.ParseResponse(res)["next"]
		if !linkOk {
			// No next pages, so break our of infinite loop here
			break
		}
		url = link.URI
	}

	// Now enrich organizations list via GraphQL
	for ix := range organizations {
		var org *Organization = &organizations[ix]
		err := client.Query(org, queryOrganizationCount, org.Slug)
		if err != nil {
			return nil, err
		}
		err = client.Query(org, queryOrganizationLists,
			org.Slug,
			org.Agents.Count,
			org.Members.Count,
			org.Pipelines.Count,
			org.SsoProviders.Count,
			org.Teams.Count)
		if err != nil {
			return nil, err
		}
	}
	return organizations, nil
}
