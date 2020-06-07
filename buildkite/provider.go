package buildkite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// SchemaMode defines the type of schema to generate for a Terraform resource.
type SchemaMode int

const (
	// DataSourceFullEntity requests the full entity schema for a Terraform resource.
	DataSourceFullEntity SchemaMode = iota
	// DataSourceReferenceOnly requests only the unique identifier for a Terraform resource.
	DataSourceReferenceOnly
)

// Provider returns the buildkite terraform provider with its schema and handlers.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// TODO -- match description to the actual permissions used
			"api_token": &schema.Schema{
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"BUILDKITE_API_TOKEN",
					"BUILDKITE_AGENT_ACCESS_TOKEN",
				}, nil),
				Description: "API token with access to GraphQL and...",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},
			"organization_slug": &schema.Schema{
				DefaultFunc: schema.EnvDefaultFunc("BUILDKITE_ORGANIZATION_SLUG", nil),
				Description: "Buildkite static organization name as used in URLs.",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},

		ConfigureFunc: providerConfigure,

		DataSourcesMap: map[string]*schema.Resource{
			"buildkite_agents":        dataSourceAgents(),
			"buildkite_builds":        dataSourceBuilds(),
			"buildkite_members":       dataSourceMembers(),
			"buildkite_organizations": dataSourceOrganizations(),
			"buildkite_pipelines":     dataSourcePipelines(),
			"buildkite_sso_providers": dataSourceSsoProviders(),
			"buildkite_teams":         dataSourceTeams(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"buildkite_pipeline":          resourcePipelineSchedule(),
			"buildkite_pipeline_schedule": resourcePipelineSchedule(),
			"buildkite_sso_provider":      resourceSsoProvider(),
			"buildkite_team":              resourceTeam(),
			"buildkite_team_member":       resourceTeamMember(),
			"buildkite_team_pipeline":     resourceTeamPipeline(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get("api_token").(string)
	organizationSlug := d.Get("organization_slug").(string)

	return NewClient(apiToken, organizationSlug), nil
}
