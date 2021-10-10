package logs

import (
	"context"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

type ContextKey string

const LogIdContextKey ContextKey = "X_LOG_ID"

const (
	DefaultConfig = `<!-- type: 设置记录器类型, 参考:https://github.com/cihub/seelog/wiki/Logger-types-reference -->
	<seelog type="asynctimer" asyncinterval="50000000" minlevel="trace" maxlevel="critical">
		<!-- <outputs> -->
		<outputs formatid="main">
			<!-- <console> -->
			<console/>
			<filter levels="trace,debug,info">
				<console formatid="colored-default"/>
			</filter>
			<filter levels="warn">
				<console formatid="colored-warn"/>
			</filter>
			<filter levels="error,critical">
				<console formatid="colored-error"/>
			</filter>
		</outputs>
		<!-- <formats>: 定制日志的输出格式, 参考:https://github.com/cihub/seelog/wiki/Format-reference -->
		<formats>
			<format id="main" format="%Level %Date(2006-01-02 15:04:05.999) %FullPath:%Line %Msg%n"/>
			<format id="colored-default"  format="%EscM(38)%Level %Date(2006-01-02 15:04:05.999) %FullPath:%Line %Msg%n%EscM(0)"/>
			<format id="colored-warn"  format="%EscM(33)%Level %Date(2006-01-02 15:04:05.999) %FullPath:%Line %Msg%n%EscM(0)"/>
			<format id="colored-error"  format="%EscM(31)%Level %Date(2006-01-02 15:04:05.999) %FullPath:%Line %Msg%n%EscM(0)"/>
		</formats>
	</seelog>`
)

var seeLogIns seelog.LoggerInterface

func init() {
	InitFromConfigAsString(DefaultConfig)
}

func InitFromConfigAsString(conf string) {
	var err error
	if seeLogIns, err = seelog.LoggerFromConfigAsBytes([]byte(conf)); err != nil {
		panic(err)
	}
	_ = seeLogIns.SetAdditionalStackDepth(1)
}

func InitFromConfigAsFile(filePath string) {
	var err error
	if seeLogIns, err = seelog.LoggerFromConfigAsFile(filePath); err != nil {
		panic(err)
	}
	_ = seeLogIns.SetAdditionalStackDepth(1)
}

func Trace(format string, params ...interface{}) {
	seeLogIns.Tracef(format, params...)
}

func Debug(format string, params ...interface{}) {
	seeLogIns.Debugf(format, params...)
}

func Info(format string, params ...interface{}) {
	seeLogIns.Infof(format, params...)
}

func Warn(format string, params ...interface{}) {
	_ = seeLogIns.Warnf(format, params...)
}

func Error(format string, params ...interface{}) {
	_ = seeLogIns.Errorf(format, params...)
}

func Critical(format string, params ...interface{}) {
	_ = seeLogIns.Criticalf(format, params...)
}

func CtxTrace(ctx context.Context, format string, params ...interface{}) {
	seeLogIns.Tracef(newFormatWithLogId(ctx, format), params...)
}

func CtxDebug(ctx context.Context, format string, params ...interface{}) {
	seeLogIns.Debugf(newFormatWithLogId(ctx, format), params...)
}

func CtxInfo(ctx context.Context, format string, params ...interface{}) {
	seeLogIns.Infof(newFormatWithLogId(ctx, format), params...)
}

func CtxWarn(ctx context.Context, format string, params ...interface{}) {
	_ = seeLogIns.Warnf(newFormatWithLogId(ctx, format), params...)
}

func CtxError(ctx context.Context, format string, params ...interface{}) {
	_ = seeLogIns.Errorf(newFormatWithLogId(ctx, format), params...)
}

func CtxCritical(ctx context.Context, format string, params ...interface{}) {
	_ = seeLogIns.Criticalf(newFormatWithLogId(ctx, format), params...)
}

func Flush() {
	seeLogIns.Flush()
}

func GenLogId() string {
	t := time.Now()
	return t.Format("20060102150405") + strconv.FormatInt(t.UnixNano()%1e9, 10)
}

func CtxWithLogId(ctx context.Context, logId string) context.Context {
	return context.WithValue(ctx, LogIdContextKey, logId)
}

func NewCtxWithLogId() context.Context {
	return CtxWithLogId(context.Background(), GenLogId())
}

func GetLogId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	logId, _ := ctx.Value(LogIdContextKey).(string)
	return logId
}

func newFormatWithLogId(ctx context.Context, format string) string {
	if ctx == nil {
		return format
	}
	if logId, _ := ctx.Value(LogIdContextKey).(string); logId != "" {
		return logId + " " + format
	}

	return format
}
