package telemetry

import (
	"os"
	"strings"
)

func IsEnabled() bool {
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "OTEL_") {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.TrimSpace(parts[1]) != "" {
			return true
		}
	}

	return false
}
