package host2

import (
	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"github.com/obscuronet/go-obscuro/go/host2/services/enclaveguardian"
	"github.com/obscuronet/go-obscuro/go/host2/services/l1"
	"github.com/obscuronet/go-obscuro/go/host2/services/l2"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

type host struct {
	// services
	l1 *l1.Service
	l2 *l2.Service

	db *db.DB
	//p2p *p2p.Service
	encl *enclaveguardian.Service

	logger gethlog.Logger
}

func NewHost(
	config *config.HostConfig,
	p2p hostcommon.P2P,
	ethClient ethadapter.EthClient,
	enclaveClient common.Enclave,
	ethWallet wallet.Wallet,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	logger gethlog.Logger,
	regMetrics gethmetrics.Registry,
) *host {

	database, err := db.CreateDBFromConfig(config, regMetrics, logger)
	if err != nil {
		logger.Crit("unable to create database for host", log.ErrKey, err)
	}

	l1Service := l1.NewL1Service(config, ethClient, ethWallet, mgmtContractLib, database, logger)
	l2Service := l2.NewL2Service(config, p2p, database, logger)
	encl := enclaveguardian.NewEnclaveGuardian(l1Service, l2Service, enclaveClient, database)
	return &host{
		l1:     l1Service,
		l2:     l2Service,
		db:     database,
		encl:   encl,
		logger: logger,
	}
}

func (h *host) Start() error {
	if err := h.l1.Start(); err != nil {
		return err
	}
	if err := h.l2.Start(); err != nil {
		return err
	}
	if err := h.encl.Start(); err != nil {
		return err
	}
	return nil
}

func (h *host) Stop() error {
	if err := h.encl.Stop(); err != nil {
		h.logger.Error("error stopping enclave guardian", log.ErrKey, err)
	}
	if err := h.l1.Stop(); err != nil {
		h.logger.Error("error stopping l1 service", log.ErrKey, err)
	}
	if err := h.l2.Stop(); err != nil {
		h.logger.Error("error stopping l2 service", log.ErrKey, err)
	}
	if err := h.db.Stop(); err != nil {
		h.logger.Error("error stopping database", log.ErrKey, err)
	}
	return nil
}
