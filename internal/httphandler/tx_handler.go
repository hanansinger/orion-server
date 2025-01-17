// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package httphandler

import (
	"net/http"
	"time"

	"github.com/hyperledger-labs/orion-server/internal/bcdb"
	internalerror "github.com/hyperledger-labs/orion-server/internal/errors"
	"github.com/hyperledger-labs/orion-server/internal/httputils"
	"github.com/hyperledger-labs/orion-server/pkg/types"
)

type txHandler struct {
	db bcdb.DB
}

// HandleTransaction handles transaction submission
func (t *txHandler) handleTransaction(w http.ResponseWriter, request *http.Request, tx interface{}, timeout time.Duration) {
	// If timeout == 0, tx is async, otherwise it is synchronous.
	resp, err := t.db.SubmitTransaction(tx, timeout)
	if err != nil {
		switch err.(type) {
		case *internalerror.BadRequestError:
			httputils.SendHTTPResponse(w, http.StatusBadRequest, &types.HttpResponseErr{ErrMsg: err.Error()})
		case *internalerror.DuplicateTxIDError:
			httputils.SendHTTPResponse(w, http.StatusBadRequest, &types.HttpResponseErr{ErrMsg: err.Error()})
		case *internalerror.TimeoutErr:
			httputils.SendHTTPResponse(w, http.StatusAccepted, &types.HttpResponseErr{ErrMsg: "Transaction processing timeout"})
		case *internalerror.NotLeaderError:
			leaderErr := err.(*internalerror.NotLeaderError)
			if leaderErr.GetLeaderID() == 0 {
				httputils.SendHTTPResponse(w, http.StatusServiceUnavailable, &types.HttpResponseErr{ErrMsg: "Cluster leader unavailable"})
			} else {
				httputils.SendHTTPRedirectServer(w, request, leaderErr.GetLeaderHostPort())
			}
		default:
			httputils.SendHTTPResponse(w, http.StatusInternalServerError, &types.HttpResponseErr{ErrMsg: err.Error()})
		}
		return
	}
	httputils.SendHTTPResponse(w, http.StatusOK, resp)
}
