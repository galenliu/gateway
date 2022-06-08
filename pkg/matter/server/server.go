package server

import (
	"github.com/galenliu/gateway/pkg/dnssd"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/matter/access"
	"github.com/galenliu/gateway/pkg/matter/config"
	"github.com/galenliu/gateway/pkg/matter/controller"
	"github.com/galenliu/gateway/pkg/matter/credentials"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/lib"
	"github.com/galenliu/gateway/pkg/matter/messageing"
	"github.com/galenliu/gateway/pkg/matter/server/internal"
	"github.com/galenliu/gateway/pkg/matter/transport"
	"log"
	"net"
	"sync"
)

type Config struct {
	UnsecureServicePort int
	SecureServicePort   int
}

type Server struct {
	mSecuredServicePort            int
	mUnsecuredServicePort          int
	mOperationalServicePort        int
	mUserDirectedCommissioningPort int
	mInterfaceId                   net.Interface
	config                         Config
	dnssdServer                    *DnssdServer
	mFabrics                       *credentials.FabricTable
	mCommissioningWindowManager    *CommissioningWindowManager
	mDeviceStorage                 lib.PersistentStorageDelegate //unknown
	mAccessControl                 access.AccessControler
	mSessionResumptionStorage      any
	mExchangeMgr                   messageing.ExchangeManager
	mAttributePersister            lib.AttributePersistenceProvider //unknown
	mAclStorage                    AclStorage
	mTransports                    transport.TransportManager
	mListener                      any
}

func NewCHIPServer() *Server {
	return &Server{}
}

func (s *Server) Init(initParams InitParams) {

	log.Printf("AppServer initializing")
	var err error

	s.mUnsecuredServicePort = initParams.OperationalServicePort
	s.mSecuredServicePort = initParams.UserDirectedCommissioningPort
	s.mInterfaceId = initParams.InterfaceId

	s.dnssdServer = NewDnssdServer()
	s.dnssdServer.SetFabricTable(s.mFabrics)

	s.mCommissioningWindowManager = NewCommissioningWindowManager(s)
	s.mCommissioningWindowManager.SetAppDelegate(initParams.AppDelegate)

	// Initialize PersistentStorageDelegate-based storage
	s.mDeviceStorage = initParams.PersistentStorageDelegate
	s.mSessionResumptionStorage = initParams.SessionResumptionStorage

	// Set up attribute persistence before we try to bring up the data model
	// handler.
	err = s.mAttributePersister.Init(s.mDeviceStorage)
	errors.SuccessOrExit(err)

	err = s.mFabrics.Init(s.mDeviceStorage)
	errors.SuccessOrExit(err)

	//少sDeviceTypeResolver参数
	err = s.mAccessControl.Init(initParams.AccessDelegate)
	errors.SuccessOrExit(err)

	s.mAclStorage = initParams.AclStorage
	err = s.mAclStorage.Init(s.mDeviceStorage, s.mFabrics)
	errors.SuccessOrExit(err)

	DnssdInstance().SetFabricTable(s.mFabrics)
	DnssdInstance().SetCommissioningModeProvider(s.mCommissioningWindowManager)

	//mGroupsProvider = initParams.groupDataProvider;
	//SetGroupDataProvider(mGroupsProvider);
	//
	//deviceInfoprovider = DeviceLayer::GetDeviceInfoProvider();
	//if (deviceInfoprovider)
	//{
	//	deviceInfoprovider->SetStorageDelegate(mDeviceStorage);
	//}

	// This initializes clusters, so should come after lower level initialization.
	//不知道干什么的
	controller.InitDataModelHandler(s.mExchangeMgr)

	params := transport.UdpListenParameters{}
	params.SetListenPort(s.mOperationalServicePort)
	params.SetNativeParams(initParams.EndpointNativeParams)
	s.mTransports, err = transport.NewUdpTransport(inet.UDPEndpointManager{}, params)
	errors.SuccessOrExit(err)

	s.mListener, err = server.IntGroupDataProviderListener(s.mTransports)
	errors.SuccessOrExit(err)

	dnssd.ResolverInstance().Init(inet.UDPEndpointManager{})

	DnssdInstance().SetSecuredPort(s.mOperationalServicePort)
	DnssdInstance().SetUnsecuredPort(s.mUserDirectedCommissioningPort)
	DnssdInstance().SetInterfaceId(s.mInterfaceId)

	if s.GetFabricTable().FabricCount() != 0 {
		if config.ConfigNetworkLayerBle {
			//TODO
			//如果Fabric不为零，那么设备已经被添加
			//可以在这里关闭蓝牙
		}
	}

	//如果设备开启了自动配对模式，进入模式
	if config.ChipDeviceConfigEnablePairingAutostart {
		s.GetFabricTable().DeleteAllFabrics()
		err = s.mCommissioningWindowManager.OpenBasicCommissioningWindow()
		errors.SuccessOrExit(err)
	}

	err = DnssdInstance().StartServer()
	errors.SuccessOrExit(err)
}

var ins *Server
var once sync.Once

// GetInstance CHIP服务，单例模式，一个应用中只会存在一个实例
func GetInstance() *Server {
	once.Do(func() {
		ins = NewCHIPServer()
	})
	return ins
}

// GetFabricTable 返回CHIP服务中的Fabric
func (s Server) GetFabricTable() *credentials.FabricTable {
	return s.mFabrics
}

func (s Server) Shutdown() {

}
