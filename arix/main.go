package main

import (
	"github.com/choria-io/fisk"
	"github.com/dyammarcano/arix/cli"
	"github.com/dyammarcano/arix/internal/version"
	"log"
	"os"
	"runtime"
)

var ver = "development"

func main() {
	help := `Arix Utility

Arix is a utility that combines various tools into a single binary.  It is

See 'arix cheat' for a quick reference to the commands available.`

	arix := fisk.New("arix", help)
	arix.Author("Arix Authors <dyam.marcano@gmail.com>")
	arix.UsageWriter(os.Stdout)
	arix.Version(version.G().Version)
	arix.WithCheats().CheatCommand.Hidden()

	opts, err := cli.ConfigureInApp(arix, nil, true)
	if err != nil {
		return
	}
	cli.SetVersion(ver)

	arix.Flag("run", "Run a command").Short('r').StringVar(&opts.Run)

	if runtime.GOOS == "windows" {

	}

	log.SetFlags(log.Ltime)
}
