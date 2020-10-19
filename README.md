# Remailer - a poor man's SRS

This project contains a small hack to get around SPF problems when forwarding
mail with a standard Unix .forward file.

If you're not familiar with the issue, you can read about SPF and SRS here:

[SRS](https://en.wikipedia.org/wiki/Sender_Rewriting_Scheme)

This project doesn't implement SRS as described above, but rather just remails
the emails, a bit like an MTU style forward.

The disadvantage of this is that bounces won't go back to the original sender
as with proper SRS, instead you'll specify an email address where all bounces
go (which will act as the sender of the remailed email).

The advantage is you can do this as a regular user, if you just have a regular
account somewhere and want to get your .forward file working again. Or if you
think the proper SRS solution is too complicated to configure for your MTA.

## How it works

After building the project (install the go toolchain and run `go build`),
assuming you're forwarding email to your account at host.com to destination1
and destination2 at example.com, you can replace a .forward file such as:

```
destination1@example.com destination2@example.com
```

with:

```
|remailer -sender sender@host.com -name "Bob Remailer" -recipients "destination1@example.com destination2@example.com" | sendmail -i -f sender@host.com "destination1@example.com destination2@example.com"
```

The remailer program will rewrite headers so it looks like the mail comes from
`sender@host.com`, but the `Reply-To` header will be set to the original sender
so you can still reply easily.

Note that the sender should probably be set to an address at the host at which
you're doing the forwarding (so the final destination will accept it), but it
shouldn't be the same account as is doing the forwarding (since then bounces
would go into a loop).
