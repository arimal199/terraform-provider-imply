package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource { return &userResource{} }

type userResource struct{ client *client.Client }

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	permissionAttrs := map[string]schema.Attribute{
		"id":        schema.StringAttribute{Computed: true},
		"name":      schema.StringAttribute{Computed: true},
		"resources": schema.ListAttribute{Computed: true, ElementType: types.StringType},
	}
	groupAttrs := map[string]schema.Attribute{
		"id":          schema.StringAttribute{Computed: true},
		"name":        schema.StringAttribute{Computed: true},
		"read_only":   schema.BoolAttribute{Computed: true},
		"permissions": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: permissionAttrs}},
		"user_count":  schema.Int64Attribute{Computed: true},
	}

	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true},
		"username":       schema.StringAttribute{Required: true},
		"email":          schema.StringAttribute{Computed: true},
		"first_name":     schema.StringAttribute{Optional: true, Computed: true},
		"last_name":      schema.StringAttribute{Optional: true, Computed: true},
		"enabled":        schema.BoolAttribute{Computed: true},
		"email_verified": schema.BoolAttribute{Computed: true},
		"permissions":    schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: permissionAttrs}},
		"groups":         schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: groupAttrs}},
		"identities":     schema.ListAttribute{Computed: true, ElementType: types.StringType},
		"actions":        schema.ListAttribute{Computed: true, ElementType: types.StringType},
		"created_on":     schema.StringAttribute{Computed: true},
	}}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	r.client = c
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{"username": plan.Username.ValueString()}
	if !plan.FirstName.IsNull() {
		body["firstName"] = plan.FirstName.ValueString()
	}
	if !plan.LastName.IsNull() {
		body["lastName"] = plan.LastName.ValueString()
	}

	_, err := r.client.Post("/users", body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	r.readIntoState(ctx, plan.Username.ValueString(), true, &plan, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.Get(fmt.Sprintf("/users/%s", state.ID.ValueString()))
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	decoded := decodeUser(user)
	if decoded.ID.IsNull() {
		decoded.ID = state.ID
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &decoded)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	_, err := r.client.Put(fmt.Sprintf("/users/%s", plan.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	user, err := r.client.Get(fmt.Sprintf("/users/%s", plan.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}
	decoded := decodeUser(user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &decoded)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(fmt.Sprintf("/users/%s", state.ID.ValueString())); err != nil && !strings.Contains(err.Error(), "status: 404") {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userResource) readIntoState(_ context.Context, username string, byUsername bool, state *UserModel, resp *resource.CreateResponse) {
	if !byUsername {
		return
	}
	users, err := r.client.Get("/users?search=" + urlQueryEscape(username))
	if err != nil {
		resp.Diagnostics.AddError("Error reading created user", err.Error())
		return
	}
	values, ok := users["values"].([]any)
	if !ok || len(values) == 0 {
		resp.Diagnostics.AddError("Error reading created user", "User was created but could not be read back")
		return
	}
	first, _ := values[0].(map[string]any)
	decoded := decodeUser(first)
	*state = decoded
}

func urlQueryEscape(v string) string {
	r := strings.NewReplacer("%", "%25", " ", "%20", "+", "%2B", "&", "%26", "=", "%3D")
	return r.Replace(v)
}
