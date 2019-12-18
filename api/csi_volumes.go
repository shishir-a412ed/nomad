package api

import (
	"sort"
	"time"
)

// CSIVolumes is used to query the top level csi volumes
type CSIVolumes struct {
	client *Client
}

// CSIVolumes returns a handle on the allocs endpoints.
func (c *Client) CSIVolumes() *CSIVolumes {
	return &CSIVolumes{client: c}
}

// List returns all CSI volumes, ignoring driver
func (v *CSIVolumes) List(q *QueryOptions) ([]*CSIVolumeListStub, *QueryMeta, error) {
	var resp []*CSIVolumeListStub
	qm, err := v.client.query("/v1/csi/volumes", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(CSIVolumeIndexSort(resp))
	return resp, qm, nil
}

// DriverList returns all CSI volumes for the specified driver
func (v *CSIVolumes) DriverList(driver string) ([]*CSIVolumeListStub, *QueryMeta, error) {
	return v.List(&QueryOptions{Prefix: driver})
}

// Info is used to retrieve a single allocation.
func (v *CSIVolumes) Info(id string, q *QueryOptions) (*CSIVolume, *QueryMeta, error) {
	var resp CSIVolume
	qm, err := v.client.query("/v1/csi/volume/"+id, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

func (v *CSIVolumes) Register(vol *CSIVolume, w *WriteOptions) error {
	req := CSIVolumeRegisterRequest{
		Volumes: []*CSIVolume{vol},
	}
	var resp struct{}
	_, err := v.client.write("/v1/csi/volume/"+vol.ID, req, &resp, w)
	return err
}

func (v *CSIVolumes) Deregister(id string, w *WriteOptions) error {
	_, err := v.client.delete("/v1/csi/volume/"+id, nil, w)
	return err
}

// CSIVolume is used for serialization
type CSIVolume struct {
	ID         string
	Driver     string
	Namespace  string
	Topology   *CSITopology
	MaxReaders int
	MaxWriters int

	// Combine structs.{Read,Write,Past}Allocs
	Allocations []*AllocationListStub

	// Healthy is true iff all the denormalized plugin health fields are true, and the
	// volume has not been marked for garbage collection
	Healthy           bool
	VolumeGC          time.Time
	ControllerID      string
	ControllerHealthy bool
	NodeHealthy       int
	NodeExpected      int

	CreatedIndex  uint64
	ModifiedIndex uint64
}

type CSIVolumeIndexSort []*CSIVolumeListStub

func (v CSIVolumeIndexSort) Len() int {
	return len(v)
}

func (v CSIVolumeIndexSort) Less(i, j int) bool {
	return v[i].CreatedIndex > v[j].CreatedIndex
}

func (v CSIVolumeIndexSort) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type CSIVolumeListStub struct {
	ID         string
	Driver     string
	Namespace  string
	Topology   *CSITopology
	MaxReaders int
	MaxWriters int
	// Healthy is true iff all the denormalized plugin health fields are true, and the
	// volume has not been marked for garbage collection
	Healthy           bool
	VolumeGC          time.Time
	ControllerID      string
	ControllerHealthy bool
	NodeHealthy       int
	NodeExpected      int
	CreatedIndex      uint64
	ModifiedIndex     uint64
}

type CSIVolumeRegisterRequest struct {
	Volumes []*CSIVolume
	WriteRequest
}

type CSIVolumeDeregisterRequest struct {
	VolumeIDs []string
	WriteRequest
}
