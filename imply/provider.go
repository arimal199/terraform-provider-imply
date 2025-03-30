package imply

import (
	"context"
	"os"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/arimal199/terraform-provider-imply/imply/polaris/auth"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &implyProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &implyProvider{
			version: version,
		}
	}
}

// implyProviderModel maps provider schema data to a Go type.
type implyProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

// implyProvider is the provider implementation.
type implyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *implyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "imply"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *implyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The Imply API host. Can be set via IMPLY_HOST environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Imply API key. Can be set via IMPLY_API_KEY environment variable.",
			},
		},
	}
}

func (p *implyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config implyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Imply API Host",
			"The provider cannot create the Imply API client as there is an unknown configuration value for the Imply API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IMPLY_HOST environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Imply API Key",
			"The provider cannot create the Imply API client as there is an unknown configuration value for the Imply API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IMPLY_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("IMPLY_HOST")
	apiKey := os.Getenv("IMPLY_API_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Imply API Host",
			"The provider cannot create the Imply API client as there is a missing or empty value for the Imply API host. "+
				"Set the host value in the configuration or use the IMPLY_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Imply API Key",
			"The provider cannot create the Imply API client as there is a missing or empty value for the Imply API key. "+
				"Set the api_key value in the configuration or use the IMPLY_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new imply client using the configuration values
	client, err := client.NewClient(&host, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Imply API Client",
			"An unexpected error occurred when creating the Imply API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Imply Client Error: "+err.Error(),
		)
		return
	}

	// Make the imply client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *implyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		auth.NewUsersDataSource,
		auth.NewGroupsDataSource,
		auth.NewPermissionsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *implyProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
