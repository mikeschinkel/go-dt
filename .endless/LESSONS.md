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

### [2026-10-30] Write release descriptions in active voice
While handing off E-2499's release tags for go-dt, I wrote release descriptions in passive or subjectless form ("Minimum Go version raised to 1.27"). Mike pointed out that this is passive voice, which he was taught is poor writing, especially in technical writing. Write release descriptions in active voice, with the release as the subject: "Raises the minimum Go version to 1.27." Whether this becomes a general writing rule is still open; Mike raised it as a possible mandate in his global Claude instructions.
- **Project**: go-dt

### [2026-10-30] Code-format literals in release descriptions
While handing off E-2499's release descriptions for go-dt, I wrote version numbers, env settings and identifiers as plain text (1.27, GOEXPERIMENT=jsonv2, reflect.Pointer). Mike wants them code-formatted in release descriptions, because the description becomes the GitHub release body, which renders Markdown: `1.27`, `GOEXPERIMENT=jsonv2`, `reflect.Pointer`.
- **Project**: go-dt

### [2026-11-05] Name a parameterized variant of a function With<Thing>, with the thing as the last parameter
While planning go-dt E-2562 (2026-11-05), Mike renamed dtx's unexported parsePathSegmentForOS(goos, s string) to the exported ParseOSPathSegmentWithOS(s, goos string). His reasons: 'With<Thing>' is the more common idiom than 'For<Thing>'; the extra thing is passed after the base function's own parameters, not before; and when a base function exists (ParseOSPathSegment), adding 'WithOS' to its name is the standard pattern he uses for the variant that takes the thing explicitly. Apply this when adding a variant of an existing function that takes one more input: keep the base name, append With<Thing>, and put the new parameter last.
- **Project**: go-dt

### [2026-11-10] go-dt types embrace constraints; never widen one to fit a new value
While planning E-2567 (EmailAddress/Mailbox in go-dt), I offered Q2 as a binary: have EmailAddress.Domain() return dt.InternetDomain after widening InternetDomain's doc to cover localhost and IP literals, or return string. Mike called it a false binary: both options violate go-dt's purpose. A go-dt type exists to carry specific constraints, so widening it weakens its value. When a value falls outside a type's constraints, propose additional constrained types (e.g. IPv4Address, IPv6Address, LocalHostname, InternetDomain) plus a broader type that any of them can be cast to (e.g. Hostname). Mike's bar for go-dt: what almost everyone would accept as correct, as if it could live in golang.org/x.
- **Project**: go-dt

### [2026-11-10] go-dt must stay backward compatible with every release; 'no callers I can see' is not evidence
While planning E-2567, I argued that tightening dt.InternetDomain was safe because nothing in go-pkgs used it. Mike replied that it breaks no callers I can see, and that go-dt should stay backward compatible with all releases unless there is a VERY STRONG reason to break. go-dt is meant for broad adoption, so unseen consumers count. Treat any change to an exported go-dt type's meaning or accepted values as a compatibility question, not one a local grep can settle.
- **Project**: go-dt

### [2026-11-10] Prefer the RFC's own term when naming go-dt types and concepts
While planning E-2567, I coined working definitions for 'label' and 'address literal' and proposed type names of my own. Mike said that if an RFC defines a term, there is a 99% chance he wants THAT term, because it is standardized. The caveats are impossible-to-resolve conflicts between RFCs, or a term also used in other relevant contexts, where he probably wants <adjective><std_term>. Also check the RFCs that updated or are referenced by the one cited (e.g. RFC 952 behind RFC 1123), not just one.
- **Project**: go-dt

### [2026-11-10] During needs-plan discussion, flesh out the space before deciding
While working through E-2567's open questions, after Mike reframed one question as a family of constrained types, I came back with seven sub-questions with recommendations plus a 'Decided without asking' list. Mike said I was jumping too quickly to locking in decisions before we had a chance to flesh out what should be decided. When Mike reopens a question's framing, explore and present the ground (terms, sources, constraints) first, and hold back recommendations and auto-decisions until he signals the space is mapped.
- **Project**: go-dt

### [2026-11-10] When Mike proposes a Go type alias, he means one constraint set under a context-specific name
While planning E-2567's type family, Mike said he could accept EmailDomain if it were an alias of InternetDomain. I answered that an alias carries no constraint of its own. Mike said I missed the point: he believed the two types had exactly the same constraints, and an alias lets different contexts use the name that reads best for them. Address the question he is asking, which is whether the constraint sets are actually equal, before explaining Go alias semantics.
- **Project**: go-dt

### [2026-11-10] In Mike's Go standards a New*() constructor never returns an error
While planning E-2567, I sketched dt.NewEmailAddress(local, host) for go-dt, which would have to fail on bad input. Mike said that in his Go standards New*() constructors should never generate an error. A constructor that can fail takes another name; Mike suggested Make*(), e.g. dt.MakeEmailAddress(EmailLocalPart, Hostname) (EmailAddress, error), and preferred deferring it because ParseEmailAddress covers the same need.
- **Project**: go-dt

### [2026-11-10] Go 1.27 has generic methods; check current Go docs before stating a language limit
While planning E-2567's HostPort type, I told Mike that Go methods can't take type parameters, so HostPort.Port() couldn't be generic. Mike replied that my training is dated and pointed to https://go.dev/blog/generic-methods: Go 1.27 adds methods with their own type parameters, e.g. func (List[E]) Map[R any](f func(E) R) List[R]. Interface methods still cannot declare type parameters, and a generic method cannot help a type satisfy an interface. go-pkgs targets go 1.27.0, so generic methods are available. Verify a claimed Go language limitation against current go.dev docs before presenting it.
- **Project**: go-dt

### [2026-11-11] Attribute a trade-off to the actual cause, and check pros against Go's conversion rules
While planning E-2567, I compared a go-dt interface (option B) with a function type (option C) for the public-suffix hook and gave B the con 'unknown TLDs still pass'. Mike asked why a method could not do what a closure can. He was right: the limit came from the method signature I had copied from cookiejar.PublicSuffixList, which returns no 'listed' flag, not from method versus closure. In the same message I claimed a typed DisplayName argument would guarantee a valid Mailbox; Mike showed NewMailbox(DisplayName("&^*&^A!@@"), ...) compiles, because any string converts to a go-dt string type. Name the real cause of each con, and do not claim a go-dt type guarantees validity, since conversion bypasses Parse.
- **Project**: go-dt

### [2026-11-11] A con shared by every option is not a con of one option
While asking how a dt.Mailbox is built for E-2567, I listed 'nothing checks a struct literal' as a con of the ParseMailbox-plus-struct-literals option only. Mike pointed out that the same con applies to every option (MakeMailbox, NewMailbox+Initialize, and Validate if it isn't called), because Mailbox's fields are exported. The needs-plan rules already exclude cons that apply equally to every option; check each con against the other options before listing it.
- **Project**: go-dt
