# Lessons Learned

A capture log of corrections, newest last. Entries are appended by
`endless lesson write`, which commits this file on the main checkout as it
writes — so recording a lesson never waits on a task landing, and never
dirties a worktree.

## Format

```
### [YYYY-MM-DD] One-line summary
<the lesson — what went wrong, why, and the rule that replaces it>
- **Project**: <project>
```

---

## Entries

### [2026-10-28] Read the rendering code before explaining why a view omits something
In go-dt ES-1369 (2026-10-27) Mike asked why E-2445's 'should follow' link to E-2444 showed no '->' in session status. I answered from the --help text that the graph draws only blocks links. Mike corrected me: '->' for precedes exists and he had seen it work. Reading internal/sessionstatuscmd/graph.go and internal/monitor/session_graph.go showed the real cause (a precedes predecessor stewarded by another session never entered the node set), filed as endless E-2455. Next time, trace the code and the stored rows before explaining a display's behaviour, rather than inferring from help text.
- **Project**: go-dt

### [2026-10-28] Agree exported API changes with Mike before making them
In go-dt E-2431 (2026-10-26) I fixed a recursion bug in dtx.ParseOSPathSegment by moving its body into an unexported helper taking the GOOS, and described it as 'refactor the dispatch logic to take the OS as a parameter'. Mike said he does not want his public APIs changed without first discussing and agreeing on it. The exported signature had not changed (apidiff clean), but the wording read as if it might. Before editing his published packages, check whether anything exported changes; if it does, propose it first; if not, say 'exported API unchanged, apidiff clean'.
- **Project**: go-dt

### [2026-10-29] Plan a task's verify suite as verify.toml; verify.sh suites are obsolete
While planning go-jsonxtractr's E-2432 (2026-10-29) I proposed a verify.sh sourcing _harness.sh, copying what E-2431 did. Mike: 'That is obsolete. Should all be verify.toml now.' Write per-task suites as verify.toml checks.
- **Project**: go-dt

### [2026-10-29] Fold adjacent work into the open task instead of listing it out of scope; ED-1550
While planning go-jsonxtractr's E-2432 (2026-10-29) I listed the repo's Makefile conversion, raising Go to 1.27, and main's uncommitted doterr.go and staged .endless/config.json as out of scope. Mike: 'Seriously? You want me file another task to do that? See ED-1550.' An exclusion with no owning task just creates a future filing. Under ED-1550, small adjacent work in the same repo goes into the task being planned; exclude only what an existing open task already owns, and name that task.
- **Project**: go-dt
