//go:build production && !dev

package bundle

import "embed"

//go:embed all:web/.output/public
var Bundle embed.FS
