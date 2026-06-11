# Bible

THe bible.db that is build is managed outside of this project, and is embedded into
the binary for distriburion. This may change later, but either way this particular
bible database will be read-only for the app.

Any "meta" stored such as highlighted sections, bookmarked regions, notes, etc will
be handled locally in on-demand-built databases (external to this) that may
reference this database's entries.

