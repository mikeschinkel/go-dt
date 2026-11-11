go-dt's EntryStatus constants (entry_status.go) and its EntryStatusFlags constant DontFollowSymlinks (entry_path_ext.go) take their values from iota. Mike's Go standards avoid iota: if such values are ever stored, for example in a database, an accidental reordering of the constants silently changes what every stored number means. These are the only iota uses in the go-dt repository, across all of its modules.

The doc comment on IsUnclassifiedEntryStatus still says it covers "a symlink, socket, device, or named pipe", though IsSymlinkEntry, IsSocketEntry, IsPipeEntry and IsDeviceEntry now classify those, and those four constants have no doc comments.

Found while planning the go-dt host, IP, DNS and email types task (E-2567, 2026-11-11). Mike asked for it as a separate task so both ship in one go-dt release rather than a later patch.
