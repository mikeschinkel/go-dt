In entry_status.go, give each constant in the EntryStatus const block an explicit value equal to its current position: IsInvalidEntryStatus 0, IsEntryError 1, IsMissingEntry 2, IsFileEntry 3, IsDirEntry 4, IsSymlinkEntry 5, IsSocketEntry 6, IsPipeEntry 7, IsDeviceEntry 8, IsUnclassifiedEntryStatus 9. Add one-line doc comments to IsSymlinkEntry, IsSocketEntry, IsPipeEntry and IsDeviceEntry, and rewrite the IsUnclassifiedEntryStatus comment to cover only entries no other constant classifies.

In entry_path_ext.go, declare DontFollowSymlinks as `1 << 0`.

Add a test to go-dt's test/ module (beside its existing tests) that asserts each EntryStatus constant and DontFollowSymlinks equals its number. Confirm with `grep -rn iota --include='*.go'` from the repo root that no iota remains.
