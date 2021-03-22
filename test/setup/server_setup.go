package setup

import (
	"crypto/tls"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
	"github.ibm.com/blockchaindb/server/internal/fileops"
	"github.ibm.com/blockchaindb/server/pkg/crypto"
	"github.ibm.com/blockchaindb/server/pkg/logger"
	"github.ibm.com/blockchaindb/server/pkg/server/mock"
	"github.ibm.com/blockchaindb/server/pkg/server/testutils"
)

// Server holds parameters related to the server
type Server struct {
	serverNum            int
	serverID             string
	address              string
	port                 int
	configDir            string
	configFilePath       string
	bootstrapFilePath    string
	cryptoMaterialsDir   string
	serverRootCACertPath string
	serverCertPath       string
	serverKeyPath        string
	adminID              string
	adminCertPath        string
	adminKeyPath         string
	adminSigner          crypto.Signer
	cmd                  *exec.Cmd
	outBuffer            *gbytes.Buffer
	errBuffer            *gbytes.Buffer
	client               *mock.Client
	logger               *logger.SugarLogger
	mu                   sync.RWMutex
}

// NewServer creates a new blockchain database server
func NewServer(id int, dir string, logger *logger.SugarLogger) (*Server, error) {
	sNumber := strconv.FormatInt(int64(id), 10)
	s := &Server{
		serverNum:          id,
		serverID:           "node-" + sNumber,
		address:            "127.0.0.1",
		port:               0,
		adminID:            "admin",
		configDir:          filepath.Join(dir, "node-"+sNumber),
		configFilePath:     filepath.Join(dir, "node-"+sNumber, "config.yml"),
		cryptoMaterialsDir: filepath.Join(dir, "node-"+sNumber, "crypto"),
		logger:             logger,
	}

	if err := fileops.CreateDir(s.configDir); err != nil {
		return nil, err
	}
	if err := fileops.CreateDir(s.cryptoMaterialsDir); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) AdminID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.adminID
}

func (s *Server) AdminSigner() crypto.Signer {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.adminSigner
}

func (s *Server) createCryptoMaterials(rootCAPemCert, caPrivKey []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyPair, err := tls.X509KeyPair(rootCAPemCert, caPrivKey)
	if err != nil {
		return err
	}
	s.serverRootCACertPath = path.Join(s.cryptoMaterialsDir, "serverRootCACert.pem")
	serverRootCACertFile, err := os.Create(s.serverRootCACertPath)
	if err != nil {
		return err
	}
	if _, err = serverRootCACertFile.Write(rootCAPemCert); err != nil {
		return err
	}
	if err = serverRootCACertFile.Close(); err != nil {
		return err
	}

	pemCert, privKey, err := testutils.IssueCertificate(s.serverID+" Instance", s.address, keyPair)
	if err != nil {
		return err
	}

	s.serverCertPath = path.Join(s.cryptoMaterialsDir, "server.pem")
	pemCertFile, err := os.Create(s.serverCertPath)
	if err != nil {
		return err
	}
	if _, err = pemCertFile.Write(pemCert); err != nil {
		return err
	}
	if err = pemCertFile.Close(); err != nil {
		return err
	}

	s.serverKeyPath = path.Join(s.cryptoMaterialsDir, "server.key")
	pemPrivKeyFile, err := os.Create(s.serverKeyPath)
	if err != nil {
		return err
	}
	if _, err = pemPrivKeyFile.Write(privKey); err != nil {
		return err
	}
	if err = pemPrivKeyFile.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) createConfigFile() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, err := os.Create(s.configFilePath)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(
		"# Integration test config.yml\n\n" +
			"node:\n" +
			"  identity:\n" +
			"    id: " + s.serverID + "\n" +
			"    certificatePath: " + s.serverCertPath + "\n" +
			"    keyPath: " + s.serverKeyPath + "\n" +
			"  network:\n" +
			"    address: " + s.address + "\n" +
			"    port: " + strconv.FormatInt(int64(s.port), 10) + "\n" +
			"  database:\n" +
			"    name: leveldb\n" +
			"    ledgerDirectory: " + filepath.Join(s.configDir, "ledger") + "\n" +
			"  queueLength:\n" +
			"    transaction: 1000\n" +
			"    reorderedTransactionBatch: 100\n" +
			"    block: 100\n" +
			"  logLevel: info\n\n" +
			"consensus:\n" +
			"  algorithm: raft\n" +
			"  maxBlockSize: 2\n" +
			"  blockTimeout: 50ms\n" +
			"  maxTransactionCountPerBlock: 1\n\n" +
			"admin:\n" +
			"  id: " + s.adminID + "\n" +
			"  certificatePath: " + s.adminCertPath + "\n\n" +
			"caconfig:\n" +
			"  rootCACertsPath: " + s.serverRootCACertPath + "\n\n",
	); err != nil {
		return err
	}

	if err = f.Sync(); err != nil {
		return err
	}

	return f.Close()
}

func (s *Server) createCmdToStartServers(executablePath string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.outBuffer = gbytes.NewBuffer()
	s.errBuffer = gbytes.NewBuffer()
	commandOut := io.MultiWriter(s.outBuffer, os.Stdout)
	commandErr := io.MultiWriter(s.errBuffer, os.Stderr)

	s.cmd = &exec.Cmd{
		Path:   executablePath,
		Args:   []string{executablePath, "start", "--configpath", s.configDir},
		Stdout: commandOut,
		Stderr: commandErr,
	}
}

func (s *Server) start(timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("Starting server " + s.serverID + " on " + s.address + ":" + strconv.FormatInt(int64(s.port), 10))
	if err := s.cmd.Start(); err != nil {
		return errors.Wrap(err, "error while starting "+s.serverID)
	}

	log.Println("Check whether the server " + s.serverID + " has started")

	g := gomega.NewWithT(&testFailure{})

	if !g.Eventually(s.outBuffer, 10).Should(gbytes.Say("Starting the server on " + s.address + ":")) {
		return errors.New("failed to start the server: " + s.serverID)
	}

	port, err := retrievePort(string(s.outBuffer.Contents()), s.address)
	if err != nil {
		return err
	}
	s.port = port

	s.logger.Debug("Successfully started server " + s.serverID + " on " + s.address + ":" + strconv.FormatInt(int64(s.port), 10))
	return nil
}

func (s *Server) shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.cmd.Process.Kill(); err != nil {
		return errors.Wrap(err, "error while shutting down "+s.serverID)
	}

	if _, err := s.cmd.Process.Wait(); err != nil {
		return errors.Wrap(err, "error while shutting down "+s.serverID)
	}

	return nil
}

func retrievePort(output string, addr string) (int, error) {
	toFind := "Starting the server on " + addr + ":"
	index := strings.Index(output, toFind)
	if index < 0 {
		return 0, errors.New("server " + addr + " has not started successfully yet")
	}

	portIndex := index + len(toFind)
	var portStr string
	for ch := output[portIndex]; ch != '\n'; ch = output[portIndex] {
		portStr += string(ch)
		portIndex++
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, err
	}
	return port, nil
}

// NewRESTClient creates a new REST client for the user to submit requests and transactions
// to the server
func (s *Server) NewRESTClient() (*mock.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return mock.NewRESTClient("http://" + s.address + ":" + strconv.FormatInt(int64(s.port), 10))
}

// testFailure is in lieu of *testing.T for gomega's types.GomegaTestingT
type testFailure struct {
}

func (t *testFailure) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
