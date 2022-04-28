package config

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/litekube/LiteKube/pkg/leader/runtime"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	"k8s.io/klog/v2"
)

type LeaderRuntime struct {
	control               *ControlSignal
	FlagsOption           *options.LeaderOptions
	RuntimeOption         *options.LeaderOptions
	RuntimeAuthentication *RuntimeAuthentications
	KineServer            *runtime.KineServer
	NetworkManagerServer  *runtime.NetWorkManager
	NetworkJoinClient     *runtime.NetWorkJoinClient
	NetworkRegisterClient *runtime.NetWorkRegisterClient
	KubernetesServer      *runtime.KubernatesServer
	controlServer         *runtime.LiteKubeControl
	OwnKineCert           bool
}

// control progress end
type ControlSignal struct {
	ctx  context.Context
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

func NewLeaderRuntime(flags *options.LeaderOptions) *LeaderRuntime {
	ctx, stop := context.WithCancel(context.TODO())
	return &LeaderRuntime{
		control: &ControlSignal{
			ctx:  ctx,
			stop: stop,
			wg:   &sync.WaitGroup{},
		},
		FlagsOption:           flags,
		RuntimeOption:         options.NewLeaderOptions(),
		RuntimeAuthentication: nil,
		NetworkManagerServer:  nil,
		NetworkJoinClient:     nil,
		NetworkRegisterClient: nil,
		KubernetesServer:      nil,
		OwnKineCert:           false,
		controlServer:         nil,
		KineServer:            nil,
	}
}

// run kine server, network manager server, network client
func (leaderRuntime *LeaderRuntime) RunForward() error {
	defer leaderRuntime.Done()
	leaderRuntime.Add()

	if leaderRuntime.RuntimeOption.GlobalOptions.RunKine {
		// run kine and network-manager
		leaderRuntime.KineServer = runtime.NewKineServer(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeOption.KineOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/kine/"),
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/kine.log"),
		)
		if err := leaderRuntime.KineServer.Run(); err != nil {
			klog.Errorf("bad args for kine server")
			return err
		}
	}

	if leaderRuntime.RuntimeOption.GlobalOptions.RunNetManager {
		leaderRuntime.NetworkManagerServer = runtime.NewNetWorkManager(leaderRuntime.control.ctx,
			leaderRuntime.RuntimeAuthentication.NetWorkManager,
			leaderRuntime.RuntimeOption.NetmamagerOptions,
			filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-manager.log"),
		)
		if err := leaderRuntime.NetworkManagerServer.Run(); err != nil {
			klog.Errorf("bad args for network manager server")
			return err
		}
	}

	leaderRuntime.NetworkJoinClient = runtime.NewNetWorkJoinClient(leaderRuntime.control.ctx,
		leaderRuntime.RuntimeOption.NetmamagerOptions,
		filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/network-client.log"),
	)
	if err := leaderRuntime.NetworkJoinClient.Run(); err != nil {
		klog.Errorf("bad args for network manager client")
		return err
	}

	leaderRuntime.NetworkRegisterClient = runtime.NewNetWorkRegisterClient(leaderRuntime.control.ctx, leaderRuntime.RuntimeOption.NetmamagerOptions)
	return nil
}

// run k8s
func (leaderRuntime *LeaderRuntime) Run() error {
	defer leaderRuntime.Done()
	leaderRuntime.Add()

	// add to same depth with LeaderRuntime.RunForward()
	// leaderRuntime.KubernetesServer = runtime.NewKubernatesServer(leaderRuntime.control.ctx,
	// 	leaderRuntime.RuntimeOption.ApiserverOptions,
	// 	leaderRuntime.RuntimeOption.ControllerManagerOptions,
	// 	leaderRuntime.RuntimeOption.SchedulerOptions,
	// 	leaderRuntime.RuntimeAuthentication.Kubernetes.KubeConfigAdmin,
	// 	filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "/logs/kubernetes/"),
	// )

	// if err := leaderRuntime.KubernetesServer.Run(); err != nil {
	// 	klog.Errorf("fail to start kubernetes server. Error: %s", err.Error())
	// 	return err
	// }

	leaderRuntime.controlServer = runtime.NewLiteKubeControl(leaderRuntime.control.ctx,
		leaderRuntime.NetworkRegisterClient,
		filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/buffer"),
		fmt.Sprintf("https://%s:%d", leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.AdvertiseAddress, leaderRuntime.RuntimeOption.ApiserverOptions.Options.SecurePort),
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.RootCaFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientCertFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.ProfessionalOptions.ClusterSigningKubeletClientKeyFile,
		leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.TokenAuthFile,
		leaderRuntime.RuntimeOption.ControllerManagerOptions.Options.ClusterCidr,
	)

	if err := leaderRuntime.controlServer.Run(); err != nil {
		klog.Errorf("fail to start litekube control server. Error: %s", err.Error())
		return err
	}

	return nil
}

func (leaderRuntime *LeaderRuntime) Stop() error {
	defer leaderRuntime.Wait()

	// give signal to end process
	leaderRuntime.control.stop()

	// stop while all return
	return nil
}

func (leaderRuntime *LeaderRuntime) Done() {
	leaderRuntime.control.wg.Done()
}

func (leaderRuntime *LeaderRuntime) Wait() {
	leaderRuntime.control.wg.Wait()
}

func (leaderRuntime *LeaderRuntime) Add() {
	leaderRuntime.control.wg.Add(1)
}
