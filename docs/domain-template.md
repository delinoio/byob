# Domain Contract Template

## Purpose
This template defines runtime and integration contracts for a project component inside a single domain.

## Required File Naming
- File name format: `docs/<domain>-<project-or-component>-<contract>.md`
- `<domain>` must match an active top-level repository domain.
- `<project-or-component>` must be lowercase kebab-case.
- `<contract>` must describe the contract purpose.

## Required Sections
All domain contract documents must include the sections below in this exact order.

## Scope
Declare the project/component and canonical implementation paths owned by this document.

## Runtime and Language
Declare runtime and primary language.

## Users and Operators
List primary users, operators, or system actors for this component.

## Interfaces and Contracts
Document public interfaces and stable identifiers.

## Storage
Document persistent data, cache, and local file contracts.

## Security
Document trust boundaries, secrets handling, authorization, and data safety constraints.

## Logging
Document required structured logs and troubleshooting expectations.

## Build and Test
Define local validation commands and CI expectations for this component.

## Dependencies and Integrations
Document upstream/downstream dependencies and cross-domain dependencies.

## Change Triggers
Declare what related docs and AGENTS contracts must be updated when this contract changes.

## References
Link to the owning `docs/project-<id>.md` index and any related contract docs.

