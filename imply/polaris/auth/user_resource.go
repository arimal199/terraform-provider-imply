package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *client.Client
}

type userResourceModel struct {
	ID            types.String      `tfsdk:"id"`
	Username      types.String      `tfsdk:"username"`
	Email         types.String      `tfsdk:"email"`
	FirstName     types.String      `tfsdk:"first_name"`
	LastName      types.String      `tfsdk:"last_name"`
	Enabled       types.Bool        `tfsdk:"enabled"`
	EmailVerified types.Bool        `tfsdk:"email_verified"`
	Permissions   []PermissionModel `tfsdk:"permissions"`
	Groups        []GroupModel      `tfsdk:"groups"`
	Actions       []types.String    `tfsdk:"actions"`
	CreatedOn     types.String      `tfsdk:"created_on"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"email_verified": schema.BoolAttribute{
				Computed: true,
			},
			"permissions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"resources": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":        schema.StringAttribute{Computed: true},
						"name":      schema.StringAttribute{Computed: true},
						"read_only": schema.BoolAttribute{Computed: true},
						"permissions": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true},
									"name": schema.StringAttribute{Computed: true},
									"resources": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"user_count": schema.Int64Attribute{Computed: true},
					},
				},
			},
			"actions": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_on": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{
		"username": plan.Username.ValueString(),
	}

	if !plan.FirstName.IsNull() {
		body["firstName"] = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		body["lastName"] = plan.LastName.ValueString()
	}
	if !plan.Enabled.IsNull() {
		body["enabled"] = plan.Enabled.ValueBool()
	}

	user, err := r.client.Post("/users", body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Imply User", err.Error())
		return
	}

	state := flattenUserResource(plan, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.Get(fmt.Sprintf("/users/%s", state.ID.ValueString()))
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read Imply User", err.Error())
		return
	}

	state = flattenUserResource(state, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	var state userResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{}
	if !plan.FirstName.IsNull() {
		body["firstName"] = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		body["lastName"] = plan.LastName.ValueString()
	}
	if !plan.Enabled.IsNull() {
		body["enabled"] = plan.Enabled.ValueBool()
	}

	user, err := r.client.Put(fmt.Sprintf("/users/%s", state.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Imply User", err.Error())
		return
	}

	nextState := flattenUserResource(plan, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &nextState)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Delete(fmt.Sprintf("/users/%s", state.ID.ValueString())); err != nil && !isNotFoundError(err) {
		resp.Diagnostics.AddError("Unable to Delete Imply User", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func flattenUserResource(plan userResourceModel, user map[string]any) userResourceModel {
	state := plan
	state.ID = stringValue(user, "id")
	state.Username = stringValue(user, "username")
	state.Email = stringValue(user, "email")
	state.FirstName = stringValue(user, "firstName")
	state.LastName = stringValue(user, "lastName")
	state.Enabled = boolValue(user, "enabled")
	state.EmailVerified = boolValue(user, "emailVerified")
	state.Permissions = permissionModels(user["permissions"])
	state.Groups = groupModels(user["groups"])
	state.Actions = stringModels(user["actions"], "")
	state.CreatedOn = stringValue(user, "createdOn")
	return state
}
