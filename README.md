# commonlib
mylog用法说明

import "github.com/kukayyou/commonlib/mylog"

初始化
mylog.InitLog(LogPath,servername, LogMaxAge, LogMaxSize, LogMaxBackups, LogLevel)
入参数说明：
LogPath,日志保存路径，如/uc/etc/
servername,server名称用于生成requestid
LogMaxAge, 日志保存时长，按天计数
LogMaxSize, 每个日志大小，单位：MB
LogMaxBackups, 日志备份个数，
LogLevel，日志级别（-1：dubug，0：info，1：warn，2：error)

使用
mylog.Debug("date now is : %s",	 time.Now().Format("2006-01-02 15:04:05"))
mylog.Info("date now is : %s",	 time.Now().Format("2006-01-02 15:04:05"))
mylog.Warn("date now is : %s",	 time.Now().Format("2006-01-02 15:04:05"))
mylog.Error("date now is : %s",	 time.Now().Format("2006-01-02 15:04:05"))
