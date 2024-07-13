package cli

import (
	"context"
	"embed"
	"github.com/choria-io/fisk"
	"github.com/dyammarcano/arix/internal/options"
	glog "log"
	"sort"
	"sync"
	"time"
)

var (
	commands []*command
	mu       sync.Mutex
	ctx      context.Context
	log      Logger

	//go:embed cheats
	fs embed.FS
)

// Logger provides a pluggable logger implementation
type Logger interface {
	Printf(format string, a ...any)
	Fatalf(format string, a ...any)
	Print(a ...any)
	Fatal(a ...any)
	Println(a ...any)
}

type goLogger struct{}

func (goLogger) Fatalf(format string, a ...any) { glog.Fatalf(format, a...) }
func (goLogger) Printf(format string, a ...any) { glog.Printf(format, a...) }
func (goLogger) Print(a ...any)                 { glog.Print(a...) }
func (goLogger) Println(a ...any)               { glog.Println(a...) }
func (goLogger) Fatal(a ...any)                 { glog.Fatal(a...) }

type command struct {
	Name    string
	Order   int
	Command func(app commandHost)
}

type commandHost interface {
	Command(name string, help string) *fisk.CmdClause
}

// SetLogger sets a custom logger to use
func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()

	log = l
}

// SetContext sets the context to use
func SetContext(c context.Context) {
	mu.Lock()
	defer mu.Unlock()

	ctx = c
}

func ConfigureInApp(app *fisk.Application, cliOpts *options.Options, prepare bool, disable ...string) (*options.Options, error) {
	if err := commonConfigure(app, cliOpts, disable...); err != nil {
		return nil, err
	}

	if prepare {
		app.PreAction(preAction)
	}

	return options.DefaultOptions, nil
}

func preAction(_ *fisk.ParseContext) (err error) {
	return nil
}

func commonConfigure(cmd commandHost, cliOpts *options.Options, disable ...string) error {
	if cliOpts != nil {
		options.DefaultOptions = cliOpts
	} else {
		options.DefaultOptions = &options.Options{
			Timeout: 5 * time.Second,
		}
	}

	if options.DefaultOptions.PrometheusNamespace == "" {
		options.DefaultOptions.PrometheusNamespace = "nats_server_check"
	}

	ctx = context.Background()
	log = goLogger{}

	sort.Slice(commands, func(i int, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	shouldEnable := func(name string) bool {
		for _, d := range disable {
			if d == name {
				return false
			}
		}

		return true
	}

	for _, c := range commands {
		if shouldEnable(c.Name) {
			c.Command(cmd)
		}
	}

	return nil
}

func SetVersion(v string) {
	mu.Lock()
	defer mu.Unlock()

	Version = v
}
