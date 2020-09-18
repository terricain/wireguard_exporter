# WireGuard Exporter

Simple wireguard exporter written in go, mainly to practice some more go. Inspired by these two wireguard exporters [[rust](https://github.com/MindFlavor/prometheus_wireguard_exporter), [go](https://github.com/mdlayher/wireguard_exporter)] as well as the main exporter structure based on the node_exporter

The exporter will add friendly names labels to the metrics if a CSV file of `publickey,name` is provided. See [here](https://github.com/terrycain/wireguard_exporter/blob/master/friendlynames.csv) for an example of the file.

## Usage

The usage closely follows that of node_exporter:

```
Usage: wireguard_exporter

Flags:
  -h, --help                                   Show context-sensitive help.
      --web.listen-address=":9586"             Address on which to expose metrics and web interface.
      --web.telemetry-path="/metrics"          Path under which to expose metrics.
      --web.disable-exporter-metrics           Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).
      --web.config=STRING                      [EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.
      --web.max-requests=2                     Maximum number of parallel scrape requests. Use 0 to disable.
      --wireguard.friendly-name-file=STRING    Path to public key to name mapping file.
      --log.level="info"                       Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format="logfmt"                    Output format of log messages. One of: [logfmt, json]
      --version                                Show application version.
```

## Exposed metrics

**wireguard_build_info** - Just standard version info

**wireguard_latest_handshake_seconds** - Unix time of the last handshake with the client

**wireguard_received_bytes_total** - Number of bytes the server has recieved from the peer.

**wireguard_sent_bytes_total** - Number of bytes the server sent to the peer.

```
# HELP wireguard_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which wireguard was built.
# TYPE wireguard_build_info gauge
wireguard_build_info{branch="HEAD",goversion="go1.14.2",revision="72c960fd6cf36d9fdce91384a991b1d7dfa164e9",version="v0.0.0-SNAPSHOT-72c960f"} 1
# HELP wireguard_latest_handshake_seconds Seconds from the last handshake
# TYPE wireguard_latest_handshake_seconds gauge
wireguard_latest_handshake_seconds{friendly_name="",interface="wg0",public_key="00mEv1wMyzVFO/Hrt++uWlziR2ZChW5hf1N6ZxrrGRw="} 1.600419993e+09
# HELP wireguard_received_bytes_total Bytes received from the peer
# TYPE wireguard_received_bytes_total counter
wireguard_received_bytes_total{friendly_name="",interface="wg0",public_key="00mEv1wMyzVFO/Hrt++uWlziR2ZChW5hf1N6ZxrrGRw="} 1.33991636e+08
# HELP wireguard_sent_bytes_total Bytes sent to the peer
# TYPE wireguard_sent_bytes_total counter
wireguard_sent_bytes_total{friendly_name="",interface="wg0",public_key="00mEv1wMyzVFO/Hrt++uWlziR2ZChW5hf1N6ZxrrGRw="} 2.587097504e+09
```

## TODO

* Add badges :D
* Add systemd unit file
* look at tests
