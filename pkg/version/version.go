package version

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Build information
// go build -ldflags "-X version.Version=x.y.z"
var (
	AppVersion string
	Revision   string
	Branch     string
	BuildUser  string
	BuildDate  string
	AppName    = "trove-fiber"
	Uptime     = time.Now()
	GoVersion  = runtime.Version()
	Platform   = runtime.GOOS + "/" + runtime.GOARCH
)

var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
  platform:         {{.platform}}
`

// Runtime returns server runtime information
func Runtime() map[string]any {
	return map[string]any{
		"app":    AppName,
		"pid":    os.Getpid(),
		"build":  Info(),
		"uptime": Uptime,
	}
}

// Info returns version and branch information
func Info() map[string]string {
	return map[string]string{
		"version":   AppVersion,
		"branch":    Branch,
		"buildUser": BuildUser,
		"goVersion": GoVersion,
	}
}

// String returns version information string.
func String() string {
	m := map[string]string{
		"program":   AppName,
		"version":   AppVersion,
		"revision":  Revision,
		"branch":    Branch,
		"buildUser": BuildUser,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
		"platform":  Platform,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

// NewCollector exports metrics about program build info
func NewCollector(app string) prometheus.Collector {
	name := strings.Replace(app, "-", "_", -1)
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: name,
			Name:      "build_info",
			Help:      fmt.Sprintf("%s build info with platform and goversion", name),
			ConstLabels: prometheus.Labels{
				"job":       name,
				"branch":    Branch,
				"version":   AppVersion,
				"revision":  Revision,
				"platform":  Platform,
				"goversion": GoVersion,
				"builduser": BuildUser,
			},
		},
		func() float64 { return 1 },
	)
}
