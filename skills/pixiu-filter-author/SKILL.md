---
name: pixiu-filter-author
description: Creates and debugs dubbo-go-pixiu HTTP/Network filters. Use when adding, registering, or wiring Pixiu filters, or when a filter compiles but does not run. Do not use for adapter or protocol-listener SPI work.
---

# pixiu-filter-author

## Purpose
Build Pixiu HTTP/Network filters that register and run.

## When to use
- Use when:
  - Adding, fixing, or wiring Pixiu HTTP/Network filters.
  - Porting Envoy/Kong/APISIX filters.
  - Diagnosing filters that compile but do not run.

- Do not use for:
  - Registry adapter or protocol-listener SPI implementation.

## Inputs

<HARD-GATE>
Do not generate YAML, write code, create files, or take any implementation action until the user has provided all required inputs. This is a first principle.

If any required input is missing, this turn must only ask for the missing fields in the current input group; do not generate examples, defaults, YAML, code, or final output.

Even if configuration information seems inferable, obvious, or implied by context, you must still ask the user to confirm it. Do not proceed until the user confirms it.
</HARD-GATE>

- Choose filter type (required):
  - Required:
    - `filter_type`: `http` or `network`
- Common plugin information (required):
  - Required:
    - `Kind` constant value, also used as YAML `name`
- HTTP filter (required when `filter_type: http`):
  - Required:
    - Which stages to implement: `Decode`, `Encode`, or both
    - Processing logic for each stage, in natural language or pseudocode
  - Optional:
    - `Config` fields and validation logic
- Network filter (required when `filter_type: network`):
  - Required:
    - Which `NetworkFilter` callbacks to implement, such as `OnData`, `OnTripleData`, or `ServeHTTP`
    - Processing logic for each callback, in natural language or pseudocode
  - Optional:
    - `Config` fields and validation logic

## Workflow
1. Check required Inputs first:
   1. Read existing code, config, and already provided user information first; do not ask again for information already present or stated.
   2. Ask one `Inputs` group at a time. Each time, output only that group's required fields, with a short explanation after each field.
   3. If the current group's required fields are incomplete, ask only for the missing fields and do not move to the next group.
   4. After the current group's required fields are complete, ask whether to fill that group's optional fields; if yes, list those optional fields with short explanations.
   5. After optional fields are skipped or completed, apply that group's defaults to complete implementation or mounting snippets; defaults must not override existing code, config, or user input.
2. Read current source before writing code:
   - `pkg/common/extension/filter/filter.go` for SPI interfaces and `pkg/pluginregistry/registry.go` for the current blank-import list.
   - HTTP filter example: `pkg/filter/cors/cors.go`.
   - Network filter example: `pkg/filter/network/dubboproxy/`.
3. Choose package shape by `filter_type`: HTTP filters go under `pkg/filter/<name>/` or `pkg/filter/http/<name>/`; Network filters go under `pkg/filter/network/<name>/`.
4. Implement the matching SPI: HTTP defines `Kind`, `Plugin`, `FilterFactory`, `Filter`, and `Config`; Network defines `Kind`, `Plugin`, `Filter`, and `Config`.
5. Finish registration and mounting: add a blank import in `pkg/pluginregistry/registry.go`, then mount according to filter type under HCM or the listener filter-chain.

## Output format
- Show the relevant YAML fragment.

## Validation
- Verify `pkg/pluginregistry/registry.go` blank-imports the new package.
- Verify HTTP filters are nested under HCM `http_filters[]`, and Network filters are under listener `filter_chains.filters[]`.
- Verify `PrepareFilterChain` deep-copies `factory.cfg` fields into each per-request Filter instance and does not share pointers.