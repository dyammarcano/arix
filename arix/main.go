package main

import (
	"github.com/choria-io/fisk"
	"github.com/dyammarcano/arix/internal/version"
	"os"
)

func main() {
	help := `Arix Utility

Arix is a utility that combines various tools into a single binary.  It is

See 'arix cheat' for a quick reference to the commands available.`

	arix := fisk.New("arix", help)
	arix.Author("Arix Authors <dyam.marcano@gmail.com>")
	arix.UsageWriter(os.Stdout)
	arix.Version(version.G().Version)
}
