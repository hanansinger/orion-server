// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hyperledger-labs/orion-server/config"
	"github.com/hyperledger-labs/orion-server/internal/httputils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestWriteLocalConfig(t *testing.T) {
	testDir, err := ioutil.TempDir("", "integration-test-setup-")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	localConfig := &config.LocalConfiguration{
		Server: config.ServerConf{
			Identity: config.IdentityConf{
				ID:              "bla",
				CertificatePath: "bla",
				KeyPath:         "bla",
			},
			Network:     config.NetworkConf{},
			Database:    config.DatabaseConf{},
			QueueLength: config.QueueLengthConf{},
			LogLevel:    "info",
		},
		BlockCreation: config.BlockCreationConf{
			MaxBlockSize:                7,
			MaxTransactionCountPerBlock: 7,
			BlockTimeout:                7,
		},
		Replication: config.ReplicationConf{
			WALDir:  "bla",
			SnapDir: "bla",
			Network: config.NetworkConf{},
			TLS: config.TLSConf{
				Enabled:               false,
				ClientAuthRequired:    false,
				ServerCertificatePath: "bla",
				ServerKeyPath:         "bla",
				ClientCertificatePath: "bla",
				ClientKeyPath:         "bla",
				CaConfig:              config.CAConfiguration{},
			},
		},
		Bootstrap: config.BootstrapConf{
			Method: "genesis",
			File:   "bla",
		},
	}

	fileName := path.Join(testDir, "config.yml")

	err = WriteLocalConfig(localConfig, fileName)
	require.NoError(t, err)

	v := viper.New()
	v.SetConfigFile(fileName)
	err = v.ReadInConfig()
	require.NoError(t, err)
	recoveredConf := &config.LocalConfiguration{}
	err = v.UnmarshalExact(recoveredConf)
	require.NoError(t, err)
	require.Equal(t, httputils.MarshalJsonOrPanic(localConfig), httputils.MarshalJsonOrPanic(recoveredConf))
}

func TestWriteSharedConfig(t *testing.T) {
	testDir, err := ioutil.TempDir("", "integration-test-setup-")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	sharedConfig := &config.SharedConfiguration{
		Nodes: []config.NodeConf{
			config.NodeConf{NodeID: "bla", Host: "bla", Port: 777, CertificatePath: "bla"},
		},
		Consensus: &config.ConsensusConf{
			Algorithm: "raft",
			Members: []*config.PeerConf{
				&config.PeerConf{NodeId: "bla", RaftId: 0, PeerHost: "bla", PeerPort: 0},
			},
			Observers: []*config.PeerConf{
				&config.PeerConf{NodeId: "bla", RaftId: 0, PeerHost: "bla", PeerPort: 0},
			},
			RaftConfig: &config.RaftConf{
				TickInterval:         "777",
				ElectionTicks:        777,
				HeartbeatTicks:       777,
				MaxInflightBlocks:    777,
				SnapshotIntervalSize: 777,
			},
		},
		CAConfig: config.CAConfiguration{},
		Admin: config.AdminConf{
			ID:              "admin",
			CertificatePath: "bla",
		},
	}

	fileName := path.Join(testDir, "shared-config.yml")

	err = WriteSharedConfig(sharedConfig, fileName)
	require.NoError(t, err)

	v := viper.New()
	v.SetConfigFile(fileName)
	err = v.ReadInConfig()
	require.NoError(t, err)
	recoveredConf := &config.SharedConfiguration{}
	err = v.UnmarshalExact(recoveredConf)
	require.NoError(t, err)
	require.Equal(t, httputils.MarshalJsonOrPanic(sharedConfig), httputils.MarshalJsonOrPanic(recoveredConf))
}
