---
version: alpha
name: 'Pixiu Admin'
description: 'A gateway operations console with a dark navigation rail and restrained mint emphasis.'
colors:
  background: '#f5f7f8'
  text: '#20252b'
  navigation: '#151a1c'
  accent: '#5eead4'
  focus: '#32a996'
  danger: '#c15e63'
typography:
  body:
    fontFamily: 'DM Sans, system-ui, sans-serif'
  technical:
    fontFamily: 'Space Mono, ui-monospace, monospace'
rounded:
  control: '6px'
  panel: '8px'
spacing:
  panel: '20px'
  form: '24px'
components:
  button: 'Primary for commit actions; secondary for validation and preview.'
  panel: 'Quiet light surface and border; dark mode uses a dark surface.'
  field: 'Visible label, hint or field error, and mint focus ring.'
---

# Pixiu Admin design context

## Overview

Pixiu Admin serves operators editing gateway configuration on desktop or a narrow laptop. Its visual reference is an operations workbench: the dark rail anchors navigation, technical values use a monospaced face, and mint marks actions and current state. The Cluster editor shows the path from draft to runtime through status, diff, preview, and history. It keeps the established product register and the existing Chinese/English locale switch.

The runtime stylesheet [src/styles.css](src/styles.css) is the canonical owner of visual values. This file records the existing choices; it does not generate CSS. The Cluster surface reuses the same palette, typography, panel, button, and field treatment. Changes to durable values require updates to the stylesheet and this file together.

## Colors and typography

The page background is `colors.background`; text is `colors.text`. Navigation uses `colors.navigation`, the active rail mark uses `colors.accent`, and keyboard focus uses `colors.focus`. Errors use `colors.danger` with text and borders, never color alone. Dark mode is owned by the `.dark` rules in `src/styles.css`.

Body copy uses DM Sans with system fallbacks for Chinese glyphs. IDs, YAML, and diff paths use Space Mono. The heading and field sizes follow the adjacent API route editor.

## Layout and shapes

The existing 248px navigation rail and page content container define the desktop frame. Cluster form fields use two columns and collapse to one below 800px. Tables scroll horizontally when narrow. Controls use a 6px radius; panels use the established rounded card treatment. Panel padding and form gaps follow the 20px and 24px rhythm already present in the route editor.

## Components and interaction

- Primary actions save or publish. Preview and validation are secondary. Delete actions are visually separated and require the shared confirmation dialog.
- Each field has a persistent label and hint; validation errors replace hints, set `aria-invalid`, and receive focus after a failed action.
- The advanced YAML mode uses the same Cluster object as the form. Parse errors keep the raw input visible and block switching to the form until corrected.
- Loading, empty, error, busy, conflict, and successful publication states use explicit text. Publication feedback says that xDS propagation follows the atomic etcd write.
- Motion follows existing controls; no decorative animation is added. Reduced motion remains governed by the global stylesheet.

## Do and don't

- Do show the stable runtime ID as system-managed information, outside editable YAML.
- Do keep draft deletion separate from published Cluster deletion.
- Don't imply a successful API publish means every running gateway has already applied xDS.
- Don't introduce a new palette or card grammar for one configuration page.
