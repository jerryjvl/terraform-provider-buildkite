package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for SSO providers.
func dataSourceSsoProviders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSsoProvidersRead,

		Schema: map[string]*schema.Schema{
			"sso_providers": schemaSsoProviderList(DataSourceFullEntity),
		},
	}
}

// Read SSO providers from the Buildkite API and convert to the Terraform schema.
func dataSourceSsoProvidersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ssoProviders, err := client.readSsoProviders()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range ssoProviders {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("sso_providers", list); err != nil {
		return fmt.Errorf("error setting sso_providers: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of SSO providers.
func schemaSsoProviderList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite SSO providers",
		Elem:        schemaSsoProvider(mode),
	}
}

// Construct a Terraform schema definition for an SSO provider.
func schemaSsoProvider(mode SchemaMode) *schema.Resource {
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

// SsoProvider defines the properties on the Buildkite API to map to Terraform.
type SsoProvider struct {
	CreatedAt                      string
	CreatedBy                      User
	DisabledAt                     string
	DisabledBy                     User
	DisabledReason                 string
	EmailDomain                    string
	EmailDomainVerificationAddress string
	EmailDomainVerifiedAt          string
	EnabledAt                      string
	EnabledBy                      User
	ID                             string
	Note                           string
	SessionDurationInHours         int
	State                          string
	TestAuthorizationRequired      bool
	Type                           string
	URL                            string
	UUID                           string
}

// SsoProviderList defines the properties on the Buildkite API to map to Terraform.
type SsoProviderList struct {
	Count int
	Edges []struct {
		Node SsoProvider
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *SsoProvider) convert(mode SchemaMode) map[string]interface{} {
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
func (source *SsoProviderList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readSsoProviders() ([]SsoProvider, error) {
	return nil, nil
}
