# Documentation Catalog

## Purpose
`docs/` is the source of truth for repository contracts.
Each project must have one project index document and one or more domain contract documents.

## Naming Rules
- Project index docs: `docs/project-<project-id>.md`
- Domain contract docs: `docs/<domain>-<project-or-component>-<contract>.md`
- Domain prefix must be one of: `cmds`, `lint`, `packages`, `docs`, or `scripts` until new top-level domains are introduced.
- Use lowercase kebab-case identifiers and stable enum-style IDs in contract sections.

## Templates
- `docs/project-template.md`
- `docs/domain-template.md`

## Project Catalog

### byob
- `docs/project-byob.md`
- `docs/cmds-byob-foundation.md`
- `docs/lint-byob-runtime-foundation.md`
- `docs/packages-byob-sdk-foundation.md`
- `docs/packages-tsgo-bridge-foundation.md`
