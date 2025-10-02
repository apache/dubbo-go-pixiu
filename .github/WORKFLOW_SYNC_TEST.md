# Workflow Sync Test

**Created**: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**Purpose**: Test automatic upstream sync workflow

## Test Information

- **Test Date**: 2025-10-02
- **Test Branch**: test/workflow-sync-20251002-1513
- **C Repository**: Similarityoung/dubbo-go-pixiu
- **B Repository**: dubbo-go-pixiu/dubbo-go-pixiu
- **A Repository**: apache/dubbo-go-pixiu

## Test Checklist

When this PR is merged into B repository, the workflow should:

- [ ] Automatically trigger on PR merge
- [ ] Rebase onto apache/dubbo-go-pixiu:develop
- [ ] Remove fork-specific workflow files:
  - [ ] .github/workflows/sync-to-upstream.yml
  - [ ] .github/workflows/SYNC_WORKFLOW_GUIDE.md
- [ ] Create PR in apache/dubbo-go-pixiu
- [ ] Comment on original PR with upstream link
- [ ] Preserve original commit authorship

## Expected Behavior

1. **Trigger**: On merge to dubbo-go-pixiu/dubbo-go-pixiu:develop
2. **Sync**: Create timestamped branch (auto-sync-YYYYMMDD-HHMMSS)
3. **Clean**: Remove workflow files automatically
4. **Create PR**: To apache/dubbo-go-pixiu:develop
5. **Notify**: Post comment with upstream PR link

## Files to Verify

After upstream PR creation, verify:

✅ **Should be present**:
- This test file (.github/WORKFLOW_SYNC_TEST.md)
- All normal code changes

❌ **Should NOT be present**:
- .github/workflows/sync-to-upstream.yml
- .github/workflows/SYNC_WORKFLOW_GUIDE.md
- .github/workflows/TESTING_GUIDE.md

## Success Criteria

- [ ] Workflow runs without errors
- [ ] Upstream PR created successfully
- [ ] Workflow files removed from upstream PR
- [ ] Original PR receives bot comment
- [ ] Commit author is preserved (not bot)

---

**Note**: This is a test file. After successful verification, it should be removed from both repositories.
