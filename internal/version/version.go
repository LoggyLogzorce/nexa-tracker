package version

import _ "embed"

const Version = "1.0.0"

//go:embed changelog.md
var ChangelogMD string
