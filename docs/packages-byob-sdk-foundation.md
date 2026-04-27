# packages-byob-sdk-foundation

## Scope
- Project/component: BYOB TypeScript SDK.
- Canonical path: `packages/byob`.

## Runtime and Language
- Runtime: Node-compatible ESM package.
- Primary language: TypeScript.

## Users and Operators
- Build-tool authors defining BYOB task graphs.
- BYOB runtime implementers consuming SDK-defined build tool definitions.

## Interfaces and Contracts
- Package name: `@delinoio/byob`.
- Public helper: `defineBuildTool(definition)`.
- Public types include build values, build inputs, build outputs, build task context, build task actions, task maps, and build tool definitions.
- The initial SDK is a type-first scaffold and does not execute build tasks.

## Storage
- The package owns no persistent storage.
- Build output is emitted to `packages/byob/dist`.

## Security
- SDK types should avoid implying access to secrets by default.
- Future runtime APIs must keep environment access explicit through context values.

## Logging
- The SDK exposes an optional task context `log` callback type.
- The initial scaffold does not implement logging behavior.

## Build and Test
- Local validation: `pnpm --filter @delinoio/byob build`.
- Type check: `pnpm --filter @delinoio/byob test`.

## Dependencies and Integrations
- Depends on TypeScript for package build validation.
- Does not depend on the Go bridge directly in the initial scaffold.

## Change Triggers
- Update `docs/project-byob.md` and this file when public SDK types or package exports change.
- Update CLI or bridge docs only when SDK changes introduce runtime integration requirements.

## References
- `docs/project-byob.md`
- `docs/domain-template.md`

