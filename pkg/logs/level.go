package logs

import (
    "github.com/ncraft-io/ncraft-go/pkg/config/reader"
)

const (
    LevelPath = "logs.level"
)

func ChangeLogLevel(value reader.Value) {
    level := value.String("")
    if len(level) > 0 && (level == "debug" || level == "info" || level == "error" ||
        level == "warn" || level == "panic" || level == "fatal") {
        SetLevel(level)
    }
}
