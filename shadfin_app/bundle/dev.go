//go:build dev && !production

package bundle

import "embed"

//go:embed all:web/.nuxt
var Bundle embed.FS
