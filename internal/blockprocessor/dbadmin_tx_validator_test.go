package blockprocessor

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.ibm.com/blockchaindb/server/internal/identity"
	"github.ibm.com/blockchaindb/server/internal/server/testutils"
	"github.ibm.com/blockchaindb/server/internal/worldstate"
	"github.ibm.com/blockchaindb/server/pkg/types"
)

func TestValidateDBAdminTx(t *testing.T) {
	t.Parallel()

	cryptoDir := testutils.GenerateTestClientCrypto(t, []string{"userWithMorePrivilege", "userWithLessPrivilege"})
	adminCert, adminSigner := testutils.LoadTestClientCrypto(t, cryptoDir, "userWithMorePrivilege")
	nonAdminCert, nonAdminSigner := testutils.LoadTestClientCrypto(t, cryptoDir, "userWithLessPrivilege")

	sampleMetadataData := &types.Metadata{
		Version: &types.Version{
			BlockNum: 2,
			TxNum:    1,
		},
	}

	userWithLessPrivilege := &types.User{
		ID:          "userWithLessPrivilege",
		Certificate: nonAdminCert.Raw,
	}
	userWithLessPrivilegeSerialized, err := proto.Marshal(userWithLessPrivilege)
	require.NoError(t, err)

	underPrivilegedUser := []*worldstate.DBUpdates{
		{
			DBName: worldstate.UsersDBName,
			Writes: []*worldstate.KVWithMetadata{
				{
					Key:      string(identity.UserNamespace) + "userWithLessPrivilege",
					Value:    userWithLessPrivilegeSerialized,
					Metadata: sampleMetadataData,
				},
			},
		},
	}

	userWithMorePrivilege := &types.User{
		ID:          "userWithMorePrivilege",
		Certificate: adminCert.Raw,
		Privilege: &types.Privilege{
			DBAdministration: true,
		},
	}
	userWithMorePrivilegeSerialized, err := proto.Marshal(userWithMorePrivilege)
	require.NoError(t, err)

	privilegedUser := []*worldstate.DBUpdates{
		{
			DBName: worldstate.UsersDBName,
			Writes: []*worldstate.KVWithMetadata{
				{
					Key:      string(identity.UserNamespace) + "userWithMorePrivilege",
					Value:    userWithMorePrivilegeSerialized,
					Metadata: sampleMetadataData,
				},
			},
		},
	}

	tests := []struct {
		name           string
		setup          func(db worldstate.DB)
		txEnv          *types.DBAdministrationTxEnvelope
		expectedResult *types.ValidationInfo
	}{
		{
			name: "invalid: signature verification failure",
			setup: func(db worldstate.DB) {
				require.NoError(t, db.Commit(privilegedUser, 1))
			},
			txEnv: testutils.SignedDBAdministrationTxEnvelope(t, nonAdminSigner, &types.DBAdministrationTx{
				UserID:    "userWithMorePrivilege",
				CreateDBs: []string{"db1", "db2"},
				DeleteDBs: []string{"db3", "db4"},
			}),
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_UNAUTHORISED,
				ReasonIfInvalid: "signature verification failed: x509: ECDSA verification failure",
			},
		},
		{
			name: "invalid: user does not have db admin privilege",
			setup: func(db worldstate.DB) {
				require.NoError(t, db.Commit(underPrivilegedUser, 1))
			},
			txEnv: testutils.SignedDBAdministrationTxEnvelope(t, nonAdminSigner,
				&types.DBAdministrationTx{
					UserID:    "userWithLessPrivilege",
					CreateDBs: []string{"db1"},
				}),
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_NO_PERMISSION,
				ReasonIfInvalid: "the user [userWithLessPrivilege] has no privilege to perform database administrative operations",
			},
		},
		{
			name: "invalid: createDBs list is invalid",
			setup: func(db worldstate.DB) {
				require.NoError(t, db.Commit(privilegedUser, 1))
			},
			txEnv: testutils.SignedDBAdministrationTxEnvelope(t, adminSigner,
				&types.DBAdministrationTx{
					UserID:    "userWithMorePrivilege",
					CreateDBs: []string{"db1", "db1"},
				}),
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db1] is duplicated in the create list",
			},
		},
		{
			name: "invalid: deleteDBs list is invalid",
			setup: func(db worldstate.DB) {
				require.NoError(t, db.Commit(privilegedUser, 1))
			},
			txEnv: testutils.SignedDBAdministrationTxEnvelope(t, adminSigner, &types.DBAdministrationTx{
				UserID:    "userWithMorePrivilege",
				DeleteDBs: []string{"db1", "db1"},
			}),
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db1] does not exist in the cluster and hence, it cannot be deleted",
			},
		},
		{
			name: "valid transaction",
			setup: func(db worldstate.DB) {
				require.NoError(t, db.Commit(privilegedUser, 1))

				createDB := []*worldstate.DBUpdates{
					{
						DBName: worldstate.DatabasesDBName,
						Writes: []*worldstate.KVWithMetadata{
							{
								Key: "db3",
							},
							{
								Key: "db4",
							},
						},
					},
				}
				require.NoError(t, db.Commit(createDB, 1))
			},
			txEnv: testutils.SignedDBAdministrationTxEnvelope(t, adminSigner, &types.DBAdministrationTx{
				UserID:    "userWithMorePrivilege",
				CreateDBs: []string{"db1", "db2"},
				DeleteDBs: []string{"db3", "db4"},
			}),
			expectedResult: &types.ValidationInfo{
				Flag: types.Flag_VALID,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := newValidatorTestEnv(t)
			defer env.cleanup()

			tt.setup(env.db)

			result, err := env.validator.dbAdminTxValidator.validate(tt.txEnv)
			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestValidateCreateDBEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(db worldstate.DB)
		toCreateDBs    []string
		expectedResult *types.ValidationInfo
	}{
		{
			name:        "invalid: dbname is empty",
			toCreateDBs: []string{""},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the name of the database to be created cannot be empty",
			},
		},
		{
			name:        "invalid: system database cannot be created",
			toCreateDBs: []string{worldstate.ConfigDBName},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [" + worldstate.ConfigDBName + "] is a system database which cannot be created as it exist by default",
			},
		},
		{
			name:        "invalid: default worldstate database cannot be created",
			toCreateDBs: []string{worldstate.DefaultDBName},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [" + worldstate.DefaultDBName + "] is the system created default database for storing states and it cannot be created as it exist by default",
			},
		},
		{
			name:        "invalid: existing database cannot be created",
			toCreateDBs: []string{"db1"},
			setup: func(db worldstate.DB) {
				createDB := []*worldstate.DBUpdates{
					{
						DBName: worldstate.DatabasesDBName,
						Writes: []*worldstate.KVWithMetadata{
							{
								Key: "db1",
							},
						},
					},
				}
				require.NoError(t, db.Commit(createDB, 1))
			},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db1] already exists in the cluster and hence, it cannot be created",
			},
		},
		{
			name:        "invalid: database is duplicated in the create list",
			toCreateDBs: []string{"db1", "db1"},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db1] is duplicated in the create list",
			},
		},
		{
			name:        "valid",
			toCreateDBs: []string{"db1", "db2"},
			expectedResult: &types.ValidationInfo{
				Flag: types.Flag_VALID,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := newValidatorTestEnv(t)
			defer env.cleanup()
			if tt.setup != nil {
				tt.setup(env.db)
			}

			result := env.validator.dbAdminTxValidator.validateCreateDBEntries(tt.toCreateDBs)
			require.True(t, proto.Equal(tt.expectedResult, result))
		})
	}
}

func TestValidateDeleteDBEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(db worldstate.DB)
		toDeleteDBs    []string
		expectedResult *types.ValidationInfo
	}{
		{
			name:        "invalid: dbname is empty",
			toDeleteDBs: []string{""},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the name of the database to be deleted cannot be empty",
			},
		},
		{
			name:        "invalid: system database cannot be deleted",
			toDeleteDBs: []string{worldstate.ConfigDBName},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [" + worldstate.ConfigDBName + "] is a system database which cannot be deleted",
			},
		},
		{
			name:        "invalid: default worldstate database cannot be deleted",
			toDeleteDBs: []string{worldstate.DefaultDBName},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [" + worldstate.DefaultDBName + "] is the system created default database to store states and it cannot be deleted",
			},
		},
		{
			name:        "invalid: non-existing database cannot be deleted",
			toDeleteDBs: []string{"db3"},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db3] does not exist in the cluster and hence, it cannot be deleted",
			},
		},
		{
			name: "invalid: database is duplicated in the delete list",
			setup: func(db worldstate.DB) {
				createDB := []*worldstate.DBUpdates{
					{
						DBName: worldstate.DatabasesDBName,
						Writes: []*worldstate.KVWithMetadata{
							{
								Key: "db1",
							},
						},
					},
				}
				require.NoError(t, db.Commit(createDB, 1))
			},
			toDeleteDBs: []string{"db1", "db1"},
			expectedResult: &types.ValidationInfo{
				Flag:            types.Flag_INVALID_INCORRECT_ENTRIES,
				ReasonIfInvalid: "the database [db1] is duplicated in the delete list",
			},
		},
		{
			name: "valid",
			setup: func(db worldstate.DB) {
				createDB := []*worldstate.DBUpdates{
					{
						DBName: worldstate.DatabasesDBName,
						Writes: []*worldstate.KVWithMetadata{
							{
								Key: "db1",
							},
							{
								Key: "db2",
							},
						},
					},
				}
				require.NoError(t, db.Commit(createDB, 1))
			},
			toDeleteDBs: []string{"db1", "db2"},
			expectedResult: &types.ValidationInfo{
				Flag: types.Flag_VALID,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := newValidatorTestEnv(t)
			defer env.cleanup()
			if tt.setup != nil {
				tt.setup(env.db)
			}

			result := env.validator.dbAdminTxValidator.validateDeleteDBEntries(tt.toDeleteDBs)
			require.True(t, proto.Equal(tt.expectedResult, result))
		})
	}
}
