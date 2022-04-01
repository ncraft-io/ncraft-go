package logs

import (
    "fmt"
    "github.com/go-kit/kit/log"
)

type kitLogger struct {
    *ZapLogger
}

func KitLogger() log.Logger {
    return &kitLogger{Logger()}
}

func NewKitLogger(zap *ZapLogger) log.Logger {
    return &kitLogger{zap}
}

// err & level should be the first key
// msg may be first or only after the level
func (l kitLogger) Log(keyvals ...interface{}) error {
    length := len(keyvals)
    secondMsgField := func(l int, key interface{}) bool { return l >= 4 && fmt.Sprint(key) == "msg" }

    if length >= 2 {
        key := fmt.Sprint(keyvals[0])
        switch key {
        case "err":
            l.ZapLogger.Errorw("", keyvals...)
        case "msg":
            l.ZapLogger.Infow(fmt.Sprint(keyvals[1]), keyvals[2:]...)
        case "level":
            level := fmt.Sprint(keyvals[1])
            switch level {
            case "debug":
                if secondMsgField(length, keyvals[2]) {
                    l.ZapLogger.Debugw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
                } else {
                    l.ZapLogger.Debugw("", keyvals[2:]...)
                }
            default:
                fallthrough
            case "info":
                if secondMsgField(length, keyvals[2]) {
                    l.ZapLogger.Infow(fmt.Sprint(keyvals[3]), keyvals[4:]...)
                } else {
                    l.ZapLogger.Infow("", keyvals[2:]...)
                }
            case "warn":
                if secondMsgField(length, keyvals[2]) {
                    l.ZapLogger.Warnw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
                } else {
                    l.ZapLogger.Warnw("", keyvals[2:]...)
                }
            case "error":
                if secondMsgField(length, keyvals[2]) {
                    l.ZapLogger.Errorw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
                } else {
                    l.ZapLogger.Errorw("", keyvals[2:]...)
                }
            }
        default:
            l.ZapLogger.Infow("", keyvals...)
        }
    }

    return nil
}
