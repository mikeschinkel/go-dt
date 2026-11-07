Shape settled with Mike in E-2547:

    type EmailAddress string
    func ParseEmailAddress(s string) (EmailAddress, error)   // net/mail.ParseAddress; rejects a display name
    func (e EmailAddress) LocalPart() string                 // "mike" in mike@newclarity.net
    func (e EmailAddress) Domain() string                    // "newclarity.net"

    type Mailbox struct {
        Name  string
        Email EmailAddress
    }
    func NewMailbox(name string, email EmailAddress) Mailbox
    func ParseMailbox(s string) (Mailbox, error)             // net/mail.ParseAddress, then ParseEmailAddress on the address
    func (m Mailbox) String() string                         // "Name <email>", unquoted, as git wants

Consumers: go-gitutil's commit ident (E-2548) and go-fsfix's RepoFixtureArgs Author and Committer with its DefaultAuthor built by NewMailbox. Tag go-dt before go-gitutil and go-fsfix take the dependency.
