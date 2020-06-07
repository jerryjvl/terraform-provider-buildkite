package buildkite

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Define a Terraform data source for members.
func dataSourceMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMembersRead,

		Schema: map[string]*schema.Schema{
			"members": schemaMemberList(DataSourceFullEntity),
		},
	}
}

// Read members from the Buildkite API and convert to the Terraform schema.
func dataSourceMembersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	members, err := client.readMembers()
	if err != nil {
		return err
	}

	list := []interface{}{}
	for _, item := range members {
		list = append(list, item.convert(DataSourceFullEntity))
	}

	if err := d.Set("members", list); err != nil {
		return fmt.Errorf("error setting members: %s", err)
	}

	d.SetId(resource.UniqueId())
	return err
}

// Construct a Terraform schema definition for a list of members.
func schemaMemberList(mode SchemaMode) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of Buildkite members",
		Elem:        schemaMember(mode),
	}
}

// Construct a Terraform schema definition for a member.
func schemaMember(mode SchemaMode) *schema.Resource {
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

// Member defines the properties on the Buildkite API to map to Terraform.
type Member struct {
	UUID string
}

type User struct {
	Avatar      struct{ URL string }
	Bot         bool
	Email       string
	HasPassword bool
	ID          string
	Name        string
	UUID        string
}

// MemberList defines the properties on the Buildkite API to map to Terraform.
type MemberList struct {
	Count int
	Edges []struct {
		Node Member
	}
}

// Convert a Buildkite API type to a Terraform structure.
func (source *Member) convert(mode SchemaMode) map[string]interface{} {
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
func (source *MemberList) convert(mode SchemaMode) (result []interface{}) {
	for _, ref := range source.Edges {
		result = append(result, ref.Node.convert(mode))
	}
	return
}

func (client *Client) readMembers() ([]Member, error) {
	return nil, nil
}
