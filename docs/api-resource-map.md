# Polaris API Resource Map

This document maps the vendored OpenAPI schema to Terraform provider constructs. It is based on `docs/openapi.json` and is intended to guide provider implementation order.

## Global APIs

| API family | Paths | Methods | Terraform fit | Current status |
| --- | --- | --- | --- | --- |
| API keys | `/v1/apikeys`, `/v1/apikeys/{id}` | `GET`, `POST`, `PUT`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| API key info | `/v1/apikeyinfo` | `GET` | Data source only | Not implemented |
| Audit events | `/v1/audit/events` | `GET` | Data source only | Not implemented |
| App name | `/v1/customizations/app-name` | `GET`, `PUT`, `DELETE` | Singleton resource + data source | Not implemented |
| Logos | `/v1/customizations/logos`, `/v1/customizations/logos/{kind}` | `GET`, `PUT`, `DELETE` | Singleton or per-kind resource + data source | Not implemented |
| Theme | `/v1/customizations/theme` | `GET`, `PUT`, `PATCH`, `DELETE` | Singleton resource + data source | Not implemented |
| Permissions | `/v1/permissions` | `GET` | Data source only | Implemented |
| Users | `/v1/users`, `/v1/users/{id}` | `GET`, `POST`, `PUT`, `DELETE` | Resource + singular/plural data sources | Data sources implemented, resource added |
| Effective permissions | `/v1/users/{id}/effectivepermissions` | `GET` | Data source only | Not implemented |
| Groups | `/v1/groups`, `/v1/groups/{id}` | `GET`, `POST`, `PUT`, `DELETE` | Resource + singular/plural data sources | Data sources implemented, resource added |
| Group members | `/v1/groups/{id}/members` | `GET`, `POST`, `DELETE` | Relationship resource + data source | Resource added, data source not implemented |
| Metrics export | `/v1/metrics/export` | `GET` | Usually out of scope for Terraform | Not implemented |
| Projects control plane | `/v1/projects`, `/v1/projects/{id}`, `/v1/project`, `/v1/project/plans` | `GET`, `POST`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |

## Project-Scoped APIs

All project-scoped resources live under `/v1/projects/{projectId}` unless otherwise noted.

| API family | Representative paths | Methods | Terraform fit | Current status |
| --- | --- | --- | --- | --- |
| Alerts | `/alerts`, `/alerts/{id}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Collections | `/collections`, `/collections/{id}` | `GET`, `POST`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Collection assets | `/collections/{id}/assets` | `POST`, `DELETE` | Relationship resource | Not implemented |
| Favorites | `/favorites`, `/favorites/{assetId}` | `GET`, `POST`, `DELETE` | Relationship resource + data source | Not implemented |
| Connections | `/connections`, `/connections/{name}` | `GET`, `POST`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Connection tests | `/connections/{name}/test` | `POST` | Action only, likely not a resource | Not implemented |
| Connection metadata | `/connectionsMeta` | `GET` | Data source only | Not implemented |
| Dashboards | `/dashboards`, `/dashboards/{id}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Dashboard pages | `/dashboards/{id}/pages`, `/pages/{pageId}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Nested resource | Not implemented |
| Dashboard tiles | `/tiles`, `/tiles/{tileId}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Nested resource | Not implemented |
| Data cubes | `/data-cubes`, `/data-cubes/{id}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Dimensions | `/data-cubes/{id}/dimensions`, `/dimensions/{dimensionId}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Nested resource | Not implemented |
| Measures | `/data-cubes/{id}/measures`, `/measures/{measureId}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Nested resource | Not implemented |
| Embedding links | `/embedding-links`, `/embedding-links/{id}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Embed keys | `/embedding-links/{id}/key`, `/v0/projects/{projectId}/pivot/api/v1/embed/key/{linkId}` | `POST`, `DELETE` | Action or sensitive subresource | Not implemented |
| Events | `/events/{connectionName}` | `POST` | Action only, poor Terraform fit | Not implemented |
| Files | `/files`, `/files/{name}` | `GET`, `POST`, `DELETE` | Resource if lifecycle is durable | Not implemented |
| Jobs | `/jobs`, `/jobs/{jobId}` | `GET`, `POST`, `PUT` | Data source first, resource only if lifecycle proves stable | Not implemented |
| Job status and logs | `/jobs/{jobId}/status`, `/logs`, `/metrics`, `/progress`, `/reset` | `GET`, `POST` | Data source or action only | Not implemented |
| Lookups | `/lookups`, `/lookups/{lookupName}` | `GET`, `POST`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Lookup aliases | `/lookups/{lookupName}/aliases` | `GET`, `PUT` | Nested resource or computed subresource | Not implemented |
| Network policy | `/network-policy` | `GET`, `PATCH` | Singleton resource + data source | Not implemented |
| Query SQL | `/query/sql`, `/query/sql/statements`, `/query/sql/statements/{queryId}` | `POST`, `GET`, `DELETE` | Not a normal Terraform resource | Not implemented |
| Reports | `/reports`, `/reports/{id}` | `GET`, `POST`, `PUT`, `PATCH`, `DELETE` | Resource + singular/plural data sources | Not implemented |
| Report evaluations | `/reports/{id}/evaluations` | `GET` | Data source only | Not implemented |
| Tables | `/tables`, `/tables/{tableName}` | `GET`, `POST`, `PUT` | Resource + singular/plural data sources | Not implemented |
| Table maintenance | `/tables/{tableName}/unusedSegments` | `GET` | Data source only | Not implemented |

## First Implementation Slice

The first implemented resource wave in this repo is:

- `imply_user`
- `imply_group`
- `imply_group_member`

That slice was chosen because it is fully present in the global schema, matches the current provider's existing read-only identity support, and can be implemented safely without introducing project-scoped configuration yet.
