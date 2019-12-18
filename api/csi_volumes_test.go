package api

import (
	"testing"

	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/require"
)

func TestCSIVolumes_CRUD(t *testing.T) {
	t.Parallel()
	c, s, root := makeACLClient(t, nil, nil)
	defer s.Stop()
	v := c.CSIVolumes()

	// Successful empty result
	vols, qm, err := v.List(nil)
	require.NoError(t, err)
	require.NotEqual(t, 0, qm.LastIndex)
	require.Equal(t, 0, len(vols))

	// Authorized QueryOpts. Use the root token to just bypass ACL details
	opts := &QueryOptions{
		Region:    "global",
		Namespace: structs.DefaultNamespace,
		AuthToken: root.SecretID,
	}

	wpts := &WriteOptions{
		Region:    "global",
		Namespace: structs.DefaultNamespace,
		AuthToken: root.SecretID,
	}

	// Register a volume
	v.Register(&CSIVolume{
		ID:         "DEADBEEF-63C7-407F-AE82-C99FBEF78FEB",
		Driver:     "minnie",
		Namespace:  structs.DefaultNamespace,
		MaxReaders: 5,
		MaxWriters: 0,
		Topology:   &CSITopology{Segments: map[string]string{"foo": "bar"}},
	}, wpts)

	// Successful result with volumes
	vols, qm, err = v.List(opts)
	require.NoError(t, err)
	require.NotEqual(t, 0, qm.LastIndex)
	require.Equal(t, 1, len(vols))

	// Successful info query
	vol, qm, err := v.Info("DEADBEEF-63C7-407F-AE82-C99FBEF78FEB", opts)
	require.NoError(t, err)
	require.Equal(t, "minnie", vol.Driver)
	require.Equal(t, 5, vol.MaxReaders)

	// Deregister the volume
	err = v.Deregister("DEADBEEF-63C7-407F-AE82-C99FBEF78FEB", wpts)
	require.NoError(t, err)

	// Successful empty result
	vols, qm, err = v.List(nil)
	require.NoError(t, err)
	require.NotEqual(t, 0, qm.LastIndex)
	require.Equal(t, 0, len(vols))

	// Failed info query
	vol, qm, err = v.Info("DEADBEEF-63C7-407F-AE82-C99FBEF78FEB", opts)
	require.Error(t, err, "missing")
}
