package runtime

import (
	"context"
	"fmt"

	"github.com/k3s-io/kine/pkg/drivers/generic"
	"github.com/k3s-io/kine/pkg/endpoint" // link to github.com/Litekube/kine, we have make some addition
	"github.com/k3s-io/kine/pkg/tls"
	"github.com/litekube/LiteKube/pkg/options/leader/kine"
	"k8s.io/klog/v2"
)

type KineServer struct {
	ctx         context.Context
	LogPath     string
	DBPath      string
	BindAddress string
	Port        uint16
	CAPath      string
	CertPath    string
	KeyPath     string
}

func NewKineServer(ctx context.Context, opt *kine.KineOptions, dbPath string, logPath string) *KineServer {
	if dbPath == "" {
		dbPath = "/"
	}

	return &KineServer{
		ctx:         ctx,
		DBPath:      dbPath,
		BindAddress: opt.BindAddress,
		Port:        opt.SecurePort,
		CAPath:      opt.CACert,
		CertPath:    opt.ServerCertFile,
		KeyPath:     opt.ServerkeyFile,
		LogPath:     logPath,
	}
}

// start run in routine and no wait
func (s *KineServer) Run() error {
	klog.Info("run network manager client")

	config := endpoint.Config{
		Listener: fmt.Sprintf("%s:%d", s.BindAddress, s.Port),
		Endpoint: fmt.Sprintf("sqlite://%s", s.DBPath),
		ServerTLSConfig: tls.Config{
			CAFile:   s.CAPath,
			CertFile: s.CertPath,
			KeyFile:  s.KeyPath,
		},
		ConnectionPoolConfig: generic.ConnectionPoolConfig{
			MaxIdle:     0,
			MaxOpen:     0,
			MaxLifetime: 0,
		},
	}

	if _, err := endpoint.Listen(s.ctx, config); err != nil {
		return err
	}

	return nil
}
