package mgrpc

import (
	"context"
	"log"
)

type ILogger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnln(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Errorln(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Fatalln(ctx context.Context, args ...interface{})
}

type Logger struct {
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[Debug] "+format, args...)
}

func (l *Logger) Debugln(ctx context.Context, args ...interface{}) {
	a := []interface{}{
		"[Debug]",
	}
	a = append(a, args...)
	log.Println(a)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *Logger) Infoln(ctx context.Context, args ...interface{}) {
	a := []interface{}{
		"[INFO]",
	}
	a = append(a, args...)
	log.Println(a)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {

}

func (l *Logger) Warnln(ctx context.Context, args ...interface{}) {

}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {

}

func (l *Logger) Errorln(ctx context.Context, args ...interface{}) {

}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {

}

func (l *Logger) Fatalln(ctx context.Context, args ...interface{}) {

}
