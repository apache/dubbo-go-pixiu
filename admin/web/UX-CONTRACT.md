# Pixiu Admin workflow contract

This records the existing configuration-screen conventions and the Cluster lifecycle added for #1082. API truth remains in `admin/API.md` and the Go handlers.

## Canonical UI Map

| Capability      | Canonical owner                            | Source of truth       | Allowed variants       | Verification                       |
| --------------- | ------------------------------------------ | --------------------- | ---------------------- | ---------------------------------- |
| Form            | Cluster schema and `ClusterEditorModel.ts` | `admin/API.md`        | Create and edit drafts | Model tests and API e2e            |
| Select/Listbox  | Native select                              | Cluster schema enum   | Discovery type         | Keyboard and popup inspection      |
| Toast           | Editor notice region                       | `src/styles.css`      | Success and error      | Live region and browser inspection |
| CRUD            | `cluster-lifecycle-api.ts`                 | `admin/API.md`        | Draft and published    | Backend tests and API e2e          |
| Table Selection | No selection in Cluster list               | Cluster page workflow | None                   | Page inspection                    |

| Capability   | Owner                                                       | Cluster behavior                                                                                                                              |
| ------------ | ----------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| Navigation   | `src/app/App.tsx` and `src/features/config/ConfigPages.tsx` | Cluster list is the section entry; editing stays in the section and Back returns to the list.                                                 |
| Table        | Existing `.route-panel` table treatment in `src/styles.css` | Published runtime records and drafts are visible together; linked records appear once. Draft publication labels use the runtime list.         |
| Form         | Cluster schema endpoint and `ClusterEditorModel.ts`         | Schema supplies discovery choices and default; invalid drafts can be saved, preview checks name availability, publish validates strictly.     |
| Select       | Native `<select>`                                           | The OS owns the popup; discovery type has few fixed options and no custom geometry requirement.                                               |
| YAML         | `ClusterEditorModel.ts`                                     | Form and YAML share one object. Invalid YAML remains visible and blocks mode changes or writes. Runtime ID is server-owned.                   |
| Confirmation | `src/components/ui/ConfirmDialog.tsx`                       | Discard, delete, and rollback require an app-owned dialog with keyboard dismissal.                                                            |
| Feedback     | Existing route editor notice pattern and field errors       | Pending actions disable duplicates; errors stay visible; success names the completed operation.                                               |
| CRUD         | `src/services/cluster-lifecycle-api.ts`                     | Saving affects only drafts; publish writes runtime; deleting a draft leaves published state; referenced published Clusters cannot be removed. |

For stale draft revisions, preserve local edits and show the server conflict; refresh before reconciling. If direct legacy CRUD changed the linked runtime Cluster, delete the stale draft and seed a new one from the current runtime Cluster. The client never retries a write with a new revision automatically. The legacy Cluster CRUD API retains its original paths and response envelope. Only authenticated Admin API routes perform lifecycle mutations.
