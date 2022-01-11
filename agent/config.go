package agent

//Config
type Config struct {
	databaseDriver           string
	databaseDSN              string
	gRPCPort                 int
	gRPCTLSCert              string
	gRPCTLSKey               string
	gRPCCACert               string
	apiPort                  int
	logPath                  string
	logLevel                 string
	heartBeatInterval        int32 // 心跳间隔,单位：s
	transExpiredMinutes      int32 // 交易完成时间,单位：min
	debug                    bool
	historyRecordExpiredDays int32
	service                  DiancanService
}

type Option interface {
	apply(*Config)
}
