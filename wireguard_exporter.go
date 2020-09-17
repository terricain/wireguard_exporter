package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/node_exporter/https"
	"github.com/terrycain/wireguard_exporter/internal/filter_handler"
	"github.com/terrycain/wireguard_exporter/internal/friendlynames"
	"net/http"
	"os"
)

var CLI struct {
	// TODO change address
	WebListenAddress string `name:"web.listen-address" default:":9586" help:"Address on which to expose metrics and web interface."`
	WebTelemetryPath string `name:"web.telemetry-path" default:"/metrics" help:"Path under which to expose metrics."`
	WebDisableExporterMetrics bool `name:"web.disable-exporter-metrics" default:"false" help:"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*)."`
	WebConfig string `name:"web.config" help:"[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication."`
	WebMaxRequests int `name:"web.max-requests" default:"2" help:"Maximum number of parallel scrape requests. Use 0 to disable."`
	WireguardFriendlyNameFile string `name:"wireguard.friendly-name-file" help:"Path to public key to name mapping file."`
	LogLevel string `name:"log.level" default:"info" enum:"debug,info,warn,error" help:"Only log messages with the given severity or above. One of: [debug, info, warn, error]"`
	LogFormat string `name:"log.format" default:"logfmt" enum:"logfmt,json" help:"Output format of log messages. One of: [logfmt, json]"`
	Version bool `name:"Show application version."`
}

func main() {
	kong.Parse(&CLI)

	if CLI.Version {
		// node_exporter, version 1.0.1 (branch: HEAD, revision: 3715be6ae899f2a9b9dbfd9c39f3e09a7bd4559f)
		//  build user:       root@1f76dbbcfa55
		//  build date:       20200616-12:44:12
		//  go version:       go1.14.4
		fmt.Print(version.Print("wireguard_exporter"))
		os.Exit(0)
	}
	logConfig := &promlog.Config{
		Level: &promlog.AllowedLevel{},
		Format: &promlog.AllowedFormat{},

	}
	logConfig.Format.Set(CLI.LogFormat)
	logConfig.Level.Set(CLI.LogLevel)
	logger := promlog.New(logConfig)

	level.Info(logger).Log("msg", "Starting wireguard_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	fctx := filter_handler.HandlerContext{
		MetricsPath:            CLI.WebTelemetryPath,
		DisableExporterMetrics: CLI.WebDisableExporterMetrics,
		MaxRequests:            CLI.WebMaxRequests,
		Logger:                 logger,
	}
	if len(CLI.WireguardFriendlyNameFile) > 0 {
		if friendlyNames, err := friendlynames.ParseFriendlyNameFile(CLI.WireguardFriendlyNameFile, logger); err != nil {
			level.Error(logger).Log("msg", "Failed to parse friendly name file", "file", CLI.WireguardFriendlyNameFile, "err", err.Error())
		} else {
			fctx.FriendlyNames = friendlyNames
		}
	}
	fctx.Load()

	http.Handle(CLI.WebTelemetryPath, fctx.MetricsHandler)
	http.HandleFunc("/", fctx.IndexFunc)

	level.Info(logger).Log("msg", "Listening on", "address", CLI.WebListenAddress)
	server := &http.Server{Addr: CLI.WebListenAddress}
	if err := https.Listen(server, CLI.WebConfig, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

}
