package friendlynames

import (
	"fmt"
	"io/ioutil"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"strings"
)

func ParseFriendlyNameFile(path string, logger log.Logger) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			level.Warn(logger).Log("msg", "Line doesnt contain publickey,friendlyname, skipping", "line", line)
			continue
		}
		result[parts[0]] = parts[1]
	}

	level.Info(logger).Log("msg", fmt.Sprintf("Loaded %d friendly names", len(result)))
	return result, nil
}
