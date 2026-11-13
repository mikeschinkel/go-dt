package dt

import (
	"strings"
	"unicode/utf8"
)

// EmailDisplayName holds the human-readable name shown with an email address,
// RFC 5322's "display-name": the "John Doe" in "John Doe <john@example.com>".
//
// The name carries the Email adjective because LDAP, Windows and profile APIs
// also use "display name", with other rules.
//
// EmailDisplayName is stored and rendered as given, unquoted, as git and
// trailers such as Co-Authored-By want it. It is not RFC 2047-encoded, and
// quotes in the input are kept as part of the name.
//
// EmailDisplayName rejects:
//   - "<" and ">". RFC 5322 allows them inside a quoted display name, but
//     unquoted output, which git stores literally, would then be ambiguous.
//   - Carriage returns, line feeds and other control characters, which would
//     break a header or a git commit's author line.
//   - Invalid UTF-8.
//
// If you need a type that accepts what EmailDisplayName rejects, define your
// own type, or open an issue at https://github.com/mikeschinkel/go-dt/issues
// to propose one.
type EmailDisplayName string

// Mailbox holds an email address and an optional display name, RFC 5322's
// "mailbox": "John Doe <john@example.com>". The field names are RFC 5322's
// terms; net/mail.Address calls the first one Name.
//
// RFC 5321 uses "mailbox" for the bare address; go-dt names that EmailAddress,
// so the two do not clash in practice.
//
// Mailbox has no New* constructor: in go-dt a New* function never returns an
// error, and with exported fields no constructor could guarantee validity
// anyway. Build one with ParseMailbox or a struct literal, and check a struct
// literal with Validate.
//
// String renders "DisplayName <address>" unquoted, or "<address>" when
// DisplayName is empty, because git's --author rejects a bare address.
type Mailbox struct {
	DisplayName EmailDisplayName
	Email       EmailAddress
}

// ParseMailboxArgs holds the options for ParseMailbox. A nil
// *ParseMailboxArgs, or a zero value, rejects a UTF-8 local part and a
// Unicode host name.
type ParseMailboxArgs struct {
	// IDNAConverter, when set, uses IDNA (Internationalized Domain Names in
	// Applications; see IDNAConverter) to convert a Unicode host name to
	// A-labels; the parsed address holds the A-label form.
	IDNAConverter IDNAConverter

	// AllowSMTPUTF8 accepts a UTF-8 local part (RFC 6531, RFC 6532).
	AllowSMTPUTF8 bool
}

// ParseEmailDisplayName parses s as an email display name. It fails with
// ErrInvalidEmailDisplayName, joined with ErrUnsupportedSyntax for an angle
// bracket; see EmailDisplayName.
func ParseEmailDisplayName(s string) (dn EmailDisplayName, err error) {
	switch {
	case s == "":
		err = NewErr(ErrInvalidEmailDisplayName, ErrEmpty)
		goto end
	case strings.ContainsAny(s, "<>"):
		err = NewErr(ErrInvalidEmailDisplayName, ErrUnsupportedSyntax, ErrAngleBracket, "type", "dt.EmailDisplayName", "input", s)
		goto end
	case !utf8.ValidString(s):
		err = NewErr(ErrInvalidEmailDisplayName, ErrInvalidUTF8, "input", s)
		goto end
	}
	for _, r := range s {
		if r < 0x20 || r == 0x7f || (r >= 0x80 && r < 0xa0) {
			err = NewErr(ErrInvalidEmailDisplayName, ErrControlCharacter, "input", s)
			goto end
		}
	}
	dn = EmailDisplayName(s)
end:
	return dn, err
}

// ParseMailbox parses s as "Display Name <address>", "<address>" or a bare
// address, leaving DisplayName empty for the last two. Space around the
// display name is trimmed. args may be nil; see ParseMailboxArgs. It fails
// with ErrInvalidMailbox joined with the error for the part that failed.
func ParseMailbox(s string, args *ParseMailboxArgs) (mb Mailbox, err error) {
	var addr, name string
	var lt int

	if args == nil {
		args = &ParseMailboxArgs{}
	}
	addr = strings.TrimSpace(s)
	if strings.HasSuffix(addr, ">") {
		lt = strings.LastIndexByte(addr, '<')
		if lt < 0 {
			err = NewErr(ErrInvalidMailbox, ErrAngleBracket, "input", s)
			goto end
		}
		name = strings.TrimSpace(addr[:lt])
		addr = addr[lt+1 : len(addr)-1]
	}
	if name != "" {
		mb.DisplayName, err = ParseEmailDisplayName(name)
		if err != nil {
			err = NewErr(ErrInvalidMailbox, "input", s, err)
			goto end
		}
	}
	mb.Email, err = ParseEmailAddress(addr, &ParseEmailAddressArgs{
		IDNAConverter: args.IDNAConverter,
		AllowSMTPUTF8: args.AllowSMTPUTF8,
	})
	if err != nil {
		err = NewErr(ErrInvalidMailbox, "input", s, err)
		goto end
	}
end:
	if err != nil {
		mb = Mailbox{}
	}
	return mb, err
}

// String renders mb as "DisplayName <address>", unquoted, or "<address>"
// when DisplayName is empty.
func (mb Mailbox) String() string {
	if mb.DisplayName == "" {
		return "<" + string(mb.Email) + ">"
	}
	return string(mb.DisplayName) + " <" + string(mb.Email) + ">"
}

// Validate returns nil when mb is a valid mailbox, or ErrInvalidMailbox
// joined with the errors ParseEmailDisplayName and ParseEmailAddress would
// return for its parts. An empty DisplayName is valid; an empty Email is not.
// Like EmailAddress.Validate, it accepts a UTF-8 local part.
func (mb Mailbox) Validate() (err error) {
	var errs []error
	if mb.DisplayName != "" {
		_, err = ParseEmailDisplayName(string(mb.DisplayName))
		errs = AppendErr(errs, err)
	}
	errs = AppendErr(errs, mb.Email.Validate())
	err = CombineErrs(errs)
	if err != nil {
		err = NewErr(ErrInvalidMailbox, mb.ErrKV(), err)
	}
	return err
}

// ErrKV returns mb's rendered text as error metadata under the key "mailbox".
func (mb Mailbox) ErrKV() ErrKV {
	return kv{k: "mailbox", v: mb.String()}
}

// ErrKV returns dn as error metadata under the key "email_display_name".
func (dn EmailDisplayName) ErrKV() ErrKV {
	return kv{k: "email_display_name", v: dn}
}

// MakeMailbox builds a Mailbox from typed parts, failing when either part is
// not valid.
//
// It ships commented out: ParseMailbox and a struct literal checked with
// Validate already cover the need. If you have a use they serve poorly, open
// a pull request at https://github.com/mikeschinkel/go-dt that enables it and
// explains the use.
//
//func MakeMailbox(dn EmailDisplayName, ea EmailAddress) (mb Mailbox, err error) {
//	mb = Mailbox{DisplayName: dn, Email: ea}
//	err = mb.Validate()
//	if err != nil {
//		mb = Mailbox{}
//	}
//	return mb, err
//}
