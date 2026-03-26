# Terraform Provider Backlog

This backlog turns the current repo state and the Polaris API mapping into an execution plan for expanding the provider from a small read-only identity surface into a scoped provider with durable resources, data sources, tests, and CI support.

## Current State

- Implemented data sources:
  - `imply_users`
  - `imply_user`
  - `imply_groups`
  - `imply_group`
  - `imply_permissions`
- Implemented resources: none
- Scope model: effectively global-only
- Client model: generic JSON transport with ad hoc host rewriting
- Test coverage: none committed

## Execution Principles

- Keep Terraform Plugin Framework resources and data sources handwritten.
- Generate API models and low-level operations from a pinned Polaris OpenAPI spec.
- Split transport by Polaris scope: `global`, `regional`, and `project`.
- Finish global identity before broad project-scoped expansion.
- Add tests before broadening the surface area.
- Treat `query/sql` and `metrics/export` as out-of-scope unless a strong Terraform use case emerges.

## Epics

### E0. Typed API Foundation

Goal: replace the current generic transport with a foundation that can scale across global, regional, and project APIs.

Primary write scope:
- `api/openapi/`
- `tools/tools.go`
- `imply/client/`
- `imply/provider.go`
- `internal/polarisapi/`

Tickets:

1. `E0-1` Vendor a pinned Polaris OpenAPI spec into `api/openapi/`.
   Acceptance criteria:
   - The repo contains the authoritative Polaris spec in a stable path.
   - The spec source and version are documented in this file or adjacent docs.
   - The build does not fetch the spec dynamically.

2. `E0-2` Add API generation tooling in `tools/tools.go` and a `make generate-api` target.
   Acceptance criteria:
   - Running `make generate-api` generates typed API code deterministically.
   - Generated code is formatted and checked in.
   - Generation failures fail locally and in CI.

3. `E0-3` Introduce a runtime package for shared auth, host normalization, retries, errors, and content types.
   Acceptance criteria:
   - One tested host normalization function handles `.app.imply.io` and `.api.imply.io`.
   - Auth header construction is centralized.
   - JSON, `multipart/form-data`, `application/merge-patch+json`, and `text/plain` are supported at the transport layer.

4. `E0-4` Introduce scoped client boundaries: `global`, `regional`, `project`.
   Acceptance criteria:
   - Transport code has explicit entry points for all three Polaris scopes.
   - `provider.Configure` can derive or construct scoped clients cleanly.
   - Existing identity reads still work against the new abstraction.

### E1. Test Harness and Safety Rails

Goal: establish a stable test harness before adding a large number of endpoints.

Primary write scope:
- `imply/client/`
- `imply/provider.go`
- `imply/internal/testserver/`
- `imply/**/testdata/`
- `imply/polaris/auth/`

Tickets:

1. `E1-1` Add transport unit tests for the client package.
   Acceptance criteria:
   - Tests cover host normalization, auth headers, methods, status handling, and invalid JSON.
   - `httptest.Server` is used for request verification.
   - Client package coverage reaches at least 90%.

2. `E1-2` Add provider configuration tests.
   Acceptance criteria:
   - `host` and `api_key` schema behavior is tested.
   - `IMPLY_HOST` and `IMPLY_API_KEY` fallbacks are tested.
   - Invalid configuration produces diagnostics rather than panics.

3. `E1-3` Extract auth response decoding into testable helpers.
   Acceptance criteria:
   - Data source files stop owning raw `map[string]any` parsing directly.
   - Decoders have success, malformed payload, empty list, and type mismatch tests.

4. `E1-4` Add schema tests for existing auth data sources.
   Acceptance criteria:
   - Attribute names and required/optional/computed semantics are tested.
   - Nested shapes for users, groups, and permissions are locked down.

### E2. Finish Global Identity

Goal: turn the current read-only identity slice into a complete, typed, Terraform-native identity surface.

Primary write scope:
- `imply/provider.go`
- `imply/polaris/auth/`
- `imply/polaris/identity/`
- `imply/client/` or generated runtime integration points

Tickets:

1. `E2-1` Add `top`, `skip`, and `search` to `imply_users`.
   Acceptance criteria:
   - Terraform schema exposes the inputs supported by the Polaris API.
   - Pagination is exercised in unit tests.
   - Existing behavior remains backward compatible.

2. `E2-2` Add `imply_effective_permissions` data source.
   Acceptance criteria:
   - The data source reads `/v1/users/{id}/effectivepermissions`.
   - Permission shapes are documented and tested.

3. `E2-3` Add `imply_group_members` data source.
   Acceptance criteria:
   - The data source reads `/v1/groups/{id}/members`.
   - `top`, `skip`, and `search` are supported if exposed by the API.

4. `E2-4` Add `imply_user` resource.
   Acceptance criteria:
   - Create, read, update, delete, and import work.
   - Read removes state on 404.
   - Acceptance coverage exists under `TF_ACC=1`.

5. `E2-5` Add `imply_group` resource.
   Acceptance criteria:
   - Create, read, update, delete, and import work.
   - Writable membership is not embedded directly into the resource.

6. `E2-6` Add `imply_group_member` relationship resource.
   Acceptance criteria:
   - Resource IDs are importable as `group_id/user_id`.
   - Add and remove operations are idempotent.
   - Read-after-write consistency is covered in acceptance tests.

### E3. Global Admin APIs

Goal: add the high-value global surfaces adjacent to identity.

Primary write scope:
- `imply/provider.go`
- `imply/polaris/projects/`
- `imply/polaris/customizations/`
- `imply/polaris/apikeys/`

Tickets:

1. `E3-1` Add `imply_api_key` resource and singular/plural data sources.
   Acceptance criteria:
   - CRUD and import behavior is tested.
   - Sensitive material is handled correctly in state.

2. `E3-2` Add `imply_project` resource and singular/plural data sources.
   Acceptance criteria:
   - Global project control-plane endpoints are supported.
   - Regional listing is modeled explicitly if required by the API split.

3. `E3-3` Add singleton customization resources for app name and theme.
   Acceptance criteria:
   - Theme patch semantics are supported.
   - Singleton import behavior is defined and tested.

4. `E3-4` Add logo management after multipart support lands.
   Acceptance criteria:
   - Upload and read flows work through the shared runtime.
   - Resource design is explicit about singleton versus per-kind ownership.

### E4. Regional and Project Surfaces

Goal: expand into durable regional and project-scoped resources once identity and global admin surfaces are stable.

Primary write scope:
- `imply/provider.go`
- `imply/polaris/projects/`
- `imply/polaris/data/`
- `imply/polaris/analytics/`
- `imply/polaris/monitoring/`

Tickets:

1. `E4-1` Add explicit provider support for regional and project scope selection.
   Acceptance criteria:
   - Provider configuration or derived scope handling is documented.
   - Scope fields that require replacement are enforced consistently.

2. `E4-2` Add `network_policy` as a singleton project resource and data source.
   Acceptance criteria:
   - Singleton semantics and import format are defined.

3. `E4-3` Add `connections`, `tables`, and `lookups`.
   Acceptance criteria:
   - Each family has resource and plural/singular data source coverage where appropriate.
   - Typed models replace placeholder model-only packages.

4. `E4-4` Add `jobs` as data-source-first support.
   Acceptance criteria:
   - The design is explicitly read-only unless the lifecycle semantics justify a resource later.

5. `E4-5` Add `dashboards`, `collections`, `data_cubes`, `reports`, `alerts`, and `embedding_links`.
   Acceptance criteria:
   - The implementation follows the same typed client and test patterns established earlier.

6. `E4-6` Evaluate `files`, `events`, `query/sql`, and `metrics/export`.
   Acceptance criteria:
   - Each family is explicitly marked as supported, data-source-only, or out-of-scope.
   - Out-of-scope decisions are documented.

### E5. Docs, CI, and Release Hygiene

Goal: remove scaffold drift and make endpoint growth safe to maintain.

Primary write scope:
- `README.md`
- `docs/`
- `examples/`
- `.github/workflows/`
- `Makefile`

Tickets:

1. `E5-1` Fix provider docs and generation metadata.
   Acceptance criteria:
   - `docs/index.md` and generated docs describe `imply`, not scaffolding placeholders.
   - Example docs reflect supported resources and data sources.

2. `E5-2` Replace placeholder examples with working Terraform configurations.
   Acceptance criteria:
   - `examples/` contains identity examples for the current and first-wave resources.
   - Examples are formatted and documented.

3. `E5-3` Align GitHub Actions with the real repo layout.
   Acceptance criteria:
   - Workflows no longer reference nonexistent scaffold paths.
   - Build, lint, generation, and test jobs run against the actual module layout.

4. `E5-4` Add CI gates for generation, unit tests, schema tests, and acceptance suites.
   Acceptance criteria:
   - PR workflows run fast feedback checks.
   - Broader acceptance suites can run on a scheduled or gated basis.

## Recommended First 10 Tickets

1. `E0-1` Vendor the Polaris OpenAPI spec.
2. `E0-2` Add generation tooling and `make generate-api`.
3. `E0-3` Add shared runtime with centralized host normalization and auth handling.
4. `E1-1` Add client transport tests.
5. `E1-2` Add provider configuration tests.
6. `E1-3` Extract auth decoders into testable helpers.
7. `E2-1` Add `top`, `skip`, and `search` to `imply_users`.
8. `E2-2` Add `imply_effective_permissions`.
9. `E2-3` Add `imply_group_members`.
10. `E2-4` Add the `imply_user` resource.

## Non-Goals for the First Wave

- SQL query execution as a Terraform-managed resource
- Metrics export as Terraform-managed state
- Broad project-scoped coverage before identity and tests are complete

