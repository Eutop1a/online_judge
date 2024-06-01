package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"online-judge/server/consts/config"
	"os"
	"path"
)

// hostname: 存储主机名。
var hostname string

// init: 初始化日志系统，设置日志级别，创建日志文件，配置日志格式，并添加一个自定义的hook来处理链路追踪数据。
func init() {
	hostname, _ = os.Hostname()

	switch config.EnvCfg.LoggerLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN", "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	}

	filePath := path.Join("/var", "log", "online-judge", "online-judge.log")
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(logTraceHook{})
	logrus.SetOutput(io.MultiWriter(f, os.Stdout))

	Logger = logrus.WithFields(logrus.Fields{
		"Tied":  config.EnvCfg.TiedLogging,
		"Host":  hostname,
		"PodIP": config.EnvCfg.PodIpAddr,
	})
}

// logTraceHook: 实现了 logrus 的Hook接口，用于将日志条目中的链路追踪信息提取并附加到日志中。
type logTraceHook struct{}

func (t logTraceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t logTraceHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	// 分布式链路追踪中的 Trace ID 和 Span ID

	// Trace ID 是一个唯一标识符，用于标识整个分布式事务或请求。
	// 它通常贯穿整个分布式系统中的多个服务和组件，用于将所有相关的操作关联起来。
	// 比如，一个用户请求可能涉及多个微服务调用，每个调用都会有自己的 Span，但它们都共享同一个 Trace ID。
	sCtx := span.SpanContext()
	if sCtx.HasTraceID() {
		entry.Data["trace_id"] = sCtx.TraceID().String()
	}

	// Span ID 是另一个唯一标识符，用于标识单个操作或处理单元。
	// 每个 Span 表示一个独立的工作单元，例如一个函数调用或一个网络请求。
	// Span 是 Trace 的组成部分，每个 Span 都有一个 Span ID，并且通常会有一个父 Span，使得 Span 之间形成一个层级结构。
	if sCtx.HasTraceID() {
		entry.Data["span_id"] = sCtx.SpanID().String()
	}

	// attribute 库用于定义和操作 OpenTelemetry 中的键值对（Key-Value pairs）。
	// 这些键值对用于描述和标识链路追踪中的特定属性或特征。它们可以被添加到 Span 中，以提供更丰富的上下文信息。
	// 例如，日志的严重性（log severity）、日志消息（log message）等都可以作为属性添加到 Span 中。
	if config.EnvCfg.LoggerWithTraceState == "enable" {
		attrs := make([]attribute.KeyValue, 0)
		logSeverityKey := attribute.Key("log.severity")
		logMessageKey := attribute.Key("log.message")
		attrs = append(attrs, logSeverityKey.String(entry.Level.String()))
		attrs = append(attrs, logMessageKey.String(entry.Message))
		for key, value := range entry.Data {
			fields := attribute.Key(fmt.Sprintf("log.fields.%s", key))
			attrs = append(attrs, fields.String(fmt.Sprintf("%v", value)))
		}
		span.AddEvent("log", trace.WithAttributes(attrs...))
		if entry.Level <= logrus.ErrorLevel {
			span.SetStatus(codes.Error, entry.Message)
		}
	}
	return nil
}

var Logger *logrus.Entry

func LogService(name string) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"Service": name,
	})
}

// SetSpanError 在Span中记录错误。
func SetSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "Internal Error")
}

// SetSpanErrorWithDesc 在Span中记录带有描述的错误。
func SetSpanErrorWithDesc(span trace.Span, err error, desc string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, desc)
}

// SetSpanWithHostname 在Span中记录主机名和Pod IP地址。
func SetSpanWithHostname(span trace.Span, hostname string) {
	span.SetAttributes(attribute.String("hostname", hostname))
	span.SetAttributes(attribute.String("podIP", config.EnvCfg.PodIpAddr))
}
