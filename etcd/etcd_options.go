package etcd

import (
	"context"
	"crypto/tls"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/coreos/etcd/clientv3"
)

//除了addr必须传入 其他的支持option设置
type Options func(*clientv3.Config)

func WithDialTimeout(DialTimeout time.Duration) Options {
	return func(config *clientv3.Config) {
		config.DialTimeout = DialTimeout
	}
}

func WithAutoSyncInterval(AutoSyncInterval time.Duration) Options {
	return func(config *clientv3.Config) {
		config.AutoSyncInterval = AutoSyncInterval
	}
}

func WithDialKeepAliveTime(DialKeepAliveTime time.Duration) Options {
	return func(config *clientv3.Config) {
		config.DialKeepAliveTime = DialKeepAliveTime
	}
}

func WithDialKeepAliveTimeout(DialKeepAliveTimeout time.Duration) Options {
	return func(config *clientv3.Config) {
		config.DialKeepAliveTimeout = DialKeepAliveTimeout
	}
}

func WithMaxCallSendMsgSize(MaxCallSendMsgSize int) Options {
	return func(config *clientv3.Config) {
		config.MaxCallSendMsgSize = MaxCallSendMsgSize
	}
}

func WithMaxCallRecvMsgSize(MaxCallRecvMsgSize int) Options {
	return func(config *clientv3.Config) {
		config.MaxCallRecvMsgSize = MaxCallRecvMsgSize
	}
}

func WithTLS(TLS *tls.Config) Options {
	return func(config *clientv3.Config) {
		config.TLS = TLS
	}
}

func WithUsername(Username string) Options {
	return func(config *clientv3.Config) {
		config.Username = Username
	}
}

func WithPassword(Password string) Options {
	return func(config *clientv3.Config) {
		config.Password = Password
	}
}

func WithRejectOldCluster(RejectOldCluster bool) Options {
	return func(config *clientv3.Config) {
		config.RejectOldCluster = RejectOldCluster
	}
}

func WithDialOptions(DialOptions []grpc.DialOption) Options {
	return func(config *clientv3.Config) {
		config.DialOptions = DialOptions
	}
}

func WithLogConfig(LogConfig *zap.Config) Options {
	return func(config *clientv3.Config) {
		config.LogConfig = LogConfig
	}
}

func WithContext(Context context.Context) Options {
	return func(config *clientv3.Config) {
		config.Context = Context
	}
}

func WithPermitWithoutStream(PermitWithoutStream bool) Options {
	return func(config *clientv3.Config) {
		config.PermitWithoutStream = PermitWithoutStream
	}
}
