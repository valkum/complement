// +build !synapse_blacklist

// Rationale for being included in Synapse's blacklist: https://github.com/matrix-org/complement/issues/38

package tests

import (
	"testing"

	"github.com/matrix-org/complement/internal/b"
	"github.com/matrix-org/complement/internal/match"
	"github.com/matrix-org/complement/internal/must"
)

// Endpoint: https://matrix.org/docs/spec/client_server/r0.6.1#post-matrix-client-r0-keys-query
// This test asserts that the server is correctly rejecting input that does not match the request format given.
// Specifically, it replaces [$device_id] with { $device_id: bool } which, if not type checked, will be processed
// like an array in Python and hence go un-noticed. In Go however it will result in a 400. The correct behaviour is
// to return a 400. Element iOS uses this erroneous format.
func TestKeysQueryWithDeviceIDAsObjectFails(t *testing.T) {
	deployment := Deploy(t, "user_query_keys", b.BlueprintAlice)
	defer deployment.Destroy(t)

	userID := "@alice:hs1"
	alice := deployment.Client(t, "hs1", userID)
	res, err := alice.DoWithAuth(t, "POST", []string{"_matrix", "client", "r0", "keys", "query"}, map[string]interface{}{
		"device_keys": map[string]interface{}{
			"@bob:hs1": map[string]bool{
				"device_id1": true,
				"device_id2": true,
			},
		},
	})
	must.NotError(t, "Failed to perform POST", err)
	must.MatchResponse(t, res, match.HTTPResponse{
		StatusCode: 400,
		JSON: []match.JSON{
			match.JSONKeyEqual("errcode", "M_BAD_JSON"),
		},
	})
}
