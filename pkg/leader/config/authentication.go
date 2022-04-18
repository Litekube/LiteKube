package config

import (
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/leader/authentication"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
)

type RuntimeAuthentications struct {
	CertDir              string
	NetWorkManagerClient *authentication.NetworkManagerClient  // nil if user provide certificate
	NetWorkManager       *authentication.NetworkAuthentication // nil if network manager not run in leader
	Kine                 *authentication.KineAuthentication    // nil if not run kine in leader or provide server certificate by user
	Kubernetes           *authentication.KubernetesAuthentication
}

func NewRuntimeAuthentication(rootCertPath string) *RuntimeAuthentications {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}
	return &RuntimeAuthentications{
		CertDir:              rootCertPath,
		Kine:                 nil,
		Kubernetes:           nil,
		NetWorkManagerClient: nil,
		NetWorkManager:       nil,
	}
}
