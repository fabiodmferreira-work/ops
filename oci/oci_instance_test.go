package oci_test

import (
	"context"
	"testing"
	"time"

	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/testutils"
	"github.com/nanovms/ops/types"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"gotest.tools/assert"
)

func TestCreateInstance(t *testing.T) {
	p, c, _, _, _ := NewProvider(t)
	ctx := testutils.NewMockContext()
	instanceName := "instance-test"
	imageName := "image-test"

	ctx.Config().RunConfig.InstanceName = instanceName
	ctx.Config().CloudConfig.ImageName = imageName

	c.EXPECT().
		ListImages(context.TODO(), core.ListImagesRequest{}).
		Return(core.ListImagesResponse{
			Items: []core.Image{
				{CompartmentId: types.StringPtr("1"), DisplayName: &imageName, TimeCreated: &common.SDKTime{Time: time.Unix(1000, 0)}, SizeInMBs: types.Int64Ptr(100000)},
			},
		}, nil)

	c.EXPECT().
		LaunchInstance(context.TODO(), core.LaunchInstanceRequest{
			LaunchInstanceDetails: core.LaunchInstanceDetails{
				AvailabilityDomain: types.StringPtr(""),
				CompartmentId:      types.StringPtr(""),
				Shape:              types.StringPtr(""),
				SourceDetails: core.InstanceSourceViaImageDetails{
					ImageId: types.StringPtr("1"),
				},
			},
		})

	err := p.CreateInstance(ctx)

	assert.NilError(t, err)
}

func TestGetInstances(t *testing.T) {
	p, c, _, _, _ := NewProvider(t)
	ctx := testutils.NewMockContext()

	c.EXPECT().
		ListInstances(context.TODO(), core.ListInstancesRequest{}).
		Return(core.ListInstancesResponse{
			Items: []core.Instance{
				{CompartmentId: types.StringPtr("1"), LifecycleState: core.InstanceLifecycleStateRunning, DisplayName: types.StringPtr("instance-1"), TimeCreated: &common.SDKTime{Time: time.Unix(1000, 0)}},
				{CompartmentId: types.StringPtr("2"), LifecycleState: core.InstanceLifecycleStateStopped, DisplayName: types.StringPtr("instance-2"), TimeCreated: &common.SDKTime{Time: time.Unix(1000, 0)}},
			},
		}, nil)

	instances, err := p.GetInstances(ctx)

	assert.NilError(t, err)

	expected := []lepton.CloudInstance{
		{ID: "1", Status: "RUNNING", Name: "instance-1", Created: "a long while ago"},
		{ID: "2", Status: "STOPPED", Name: "instance-2", Created: "a long while ago"},
	}

	assert.DeepEqual(t, expected, instances)
}
