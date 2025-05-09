package logger

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type streamWrapper struct {
    l Logger
    grpc.ServerStream
}

type Logger interface {
    Info(args ...interface{})
    Warn(args ...interface{})
    Debug(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
}

type customTextFormatter struct {}

func (f *customTextFormatter) Format(l *logrus.Entry) ([]byte, error){
    var buffer bytes.Buffer

    switch l.Data["event"] {
    case "grpcRequestStart":
        buffer.WriteString(fmt.Sprintf("[%s] [%s] %s %s %v",
            strings.ToUpper(l.Level.String()),
            l.Data["method"],
            l.Message,
            l.Data["time"].(time.Time).Format("2006.01.02 15:04:05"),
            l.Data["request"],
        ))
    case "grpcRequestEnd":
        buffer.WriteString(fmt.Sprintf("[%s] [%s] [%s] Duration: %s Error: %v Status:%s Response %v",
           strings.ToUpper(l.Level.String()),
           l.Data["method"],
           l.Message,
           l.Data["duration"],
           l.Data["error"],
           l.Data["status"],
           l.Data["response"],
        ))
    default:
        buffer.WriteString(fmt.Sprintf("[%s] %s [%s] [%s:%d] %s",
            strings.ToUpper(l.Level.String()),
            time.Now().Format("2006.01.02 15:04:05"),
            l.Caller.File,
            l.Caller.Func.Name(),
            l.Caller.Line,
            l.Message,
        ))

    }
    
    buffer.WriteString("\n")
    return buffer.Bytes(), nil
}

func LoggingInterceptor(l *logrus.Logger) grpc.UnaryServerInterceptor {
    return func (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        start := time.Now()
        
        l.WithFields(logrus.Fields{
            "event":    "grpcRequestStart",
            "method":   info.FullMethod,
            "time":     start,
            "request":  req,
        }).Info("Request has started")

        resp, err := handler(ctx, req)

        end := time.Since(start)

        status, _ := status.FromError(err)

        l.WithFields(logrus.Fields{
            "event":    "grpcRequestEnd",
            "method":   info.FullMethod,
            "duration": end,
            "error":    err,
            "status":   status, 
            "response": resp,
        }).Info("Request has done")

        return resp, err
    }
}

func (w *streamWrapper) SendMsg(m interface{}) error {
    w.l.Info("Sending message:", m)
    return w.ServerStream.SendMsg(m)
}


func StreamLoggingInterceptor(l *logrus.Logger) grpc.StreamServerInterceptor {
    return func (srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        wrap := newStreamWrapper(stream, l)
        
        start := time.Now()
        l.WithFields(logrus.Fields {
            "event":   "grpcRequestStart",
            "method":  info.FullMethod,
            "time":    start,
            "request": srv,
        }).Info("Stream started")

        err := handler(srv, wrap)

        return err
    }
}

func NewLogger() *logrus.Logger {
    log := logrus.New()

    log.SetLevel(logrus.InfoLevel)
    log.SetOutput(os.Stdout)
    log.SetReportCaller(true)
    log.SetFormatter(&customTextFormatter{})


    defer log.Info("Logger has been init")
    return log
}

func newStreamWrapper(stream grpc.ServerStream, l *logrus.Logger) *streamWrapper {
    return &streamWrapper {
        l: l,
        ServerStream: stream,
    }
}
