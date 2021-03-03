package oci_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/testutils"
	"github.com/nanovms/ops/types"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"github.com/spf13/afero"
	"gotest.tools/assert"
)

type putObjectMatcher struct {
	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// The name of the object. Avoid entering confidential information.
	// Example: `test/object1.log`
	ObjectName *string `mandatory:"true" contributesTo:"path" name:"objectName"`

	// The content length of the body.
	ContentLength *int64 `mandatory:"true" contributesTo:"header" name:"Content-Length"`
}

func (p *putObjectMatcher) Matches(x interface{}) bool {
	want := x.(objectstorage.PutObjectRequest)
	return *want.BucketName == *p.BucketName && *want.NamespaceName == *p.NamespaceName && *want.ObjectName == *p.ObjectName && *want.ContentLength == *p.ContentLength
}

func (p *putObjectMatcher) String() string {
	var bucketName string
	var namespaceName string
	var objectName string
	var contentLength int

	if p.BucketName != nil {
		bucketName = *p.BucketName
	}

	if p.NamespaceName != nil {
		namespaceName = *p.NamespaceName
	}

	if p.ObjectName != nil {
		objectName = *p.ObjectName
	}

	if p.ContentLength != nil {
		contentLength = int(*p.ContentLength)
	}

	return fmt.Sprintf("{ NamespaceName=%s BucketName=%s ObjectName=%s ContentLength=%d }", namespaceName, bucketName, objectName, contentLength)
}

func PutObjectMatcher(x objectstorage.PutObjectRequest) gomock.Matcher {
	return &putObjectMatcher{
		BucketName:    x.BucketName,
		ContentLength: x.ContentLength,
		NamespaceName: x.NamespaceName,
		ObjectName:    x.ObjectName,
	}
}

func TestCreateImage(t *testing.T) {
	p, computeService, _, storageService, fs := NewProvider(t)
	image, _ := afero.TempFile(fs, "", "oci-image")
	imagePath := image.Name()
	cloudImageName := "main"
	bucketName := "test-bucket"
	bucketNamespace := "test-namespace"
	ctx := testutils.NewMockContext()

	ctx.Config().CloudConfig.ImageName = cloudImageName
	ctx.Config().CloudConfig.BucketName = bucketName
	ctx.Config().CloudConfig.BucketNamespace = bucketNamespace

	storageService.EXPECT().
		PutObject(context.TODO(), PutObjectMatcher(objectstorage.PutObjectRequest{NamespaceName: &bucketNamespace, ObjectName: &cloudImageName, BucketName: &bucketName, ContentLength: types.Int64Ptr(0)})).
		Return(objectstorage.PutObjectResponse{}, nil)

	computeService.EXPECT().
		CreateImage(context.TODO(), core.CreateImageRequest{
			CreateImageDetails: core.CreateImageDetails{
				ImageSourceDetails: core.ImageSourceViaObjectStorageTupleDetails{
					NamespaceName:   &bucketNamespace,
					BucketName:      &bucketName,
					ObjectName:      &cloudImageName,
					SourceImageType: core.ImageSourceDetailsSourceImageTypeQcow2,
				},
			},
		}).
		Return(core.CreateImageResponse{}, nil)

	err := p.CreateImage(ctx, imagePath)

	assert.NilError(t, err)
}

func TestListImages(t *testing.T) {
	p, computeService, _, _, _ := NewProvider(t)
	ctx := testutils.NewMockContext()

	computeService.EXPECT().
		ListImages(context.TODO(), core.ListImagesRequest{}).
		Return(core.ListImagesResponse{
			Items: []core.Image{
				{
					CompartmentId:  types.StringPtr("1"),
					DisplayName:    types.StringPtr("img1"),
					LifecycleState: core.ImageLifecycleStateAvailable,
					TimeCreated:    &common.SDKTime{Time: time.Unix(100000, 0)},
					SizeInMBs:      types.Int64Ptr(2000),
				},
				{
					CompartmentId:  types.StringPtr("2"),
					DisplayName:    types.StringPtr("img2"),
					LifecycleState: core.ImageLifecycleStateImporting,
					TimeCreated:    &common.SDKTime{Time: time.Unix(200000, 0)},
					SizeInMBs:      types.Int64Ptr(2000),
				},
				{
					CompartmentId:  types.StringPtr("3"),
					DisplayName:    types.StringPtr("img3"),
					LifecycleState: core.ImageLifecycleStateDeleted,
					TimeCreated:    &common.SDKTime{Time: time.Unix(300000, 0)},
					SizeInMBs:      types.Int64Ptr(3000),
				},
			},
		}, nil)

	images, err := p.GetImages(ctx)

	assert.NilError(t, err)

	assert.DeepEqual(t, []lepton.CloudImage{
		{ID: "1", Name: "img1", Status: string(core.ImageLifecycleStateAvailable), Created: time.Unix(100000, 0), Size: 1},
		{ID: "2", Name: "img2", Status: string(core.ImageLifecycleStateImporting), Created: time.Unix(200000, 0), Size: 1},
		{ID: "3", Name: "img3", Status: string(core.ImageLifecycleStateDeleted), Created: time.Unix(300000, 0), Size: 2},
	}, images)
}

func TestDeleteImage(t *testing.T) {
	p, computeService, _, _, _ := NewProvider(t)
	ctx := testutils.NewMockContext()

	computeService.EXPECT().
		ListImages(context.TODO(), core.ListImagesRequest{}).
		Return(core.ListImagesResponse{
			Items: []core.Image{
				{
					CompartmentId:  types.StringPtr("1"),
					DisplayName:    types.StringPtr("img1"),
					LifecycleState: core.ImageLifecycleStateAvailable,
					TimeCreated:    &common.SDKTime{Time: time.Unix(100000, 0)},
					SizeInMBs:      types.Int64Ptr(2000),
				},
				{
					CompartmentId:  types.StringPtr("2"),
					DisplayName:    types.StringPtr("test"),
					LifecycleState: core.ImageLifecycleStateImporting,
					TimeCreated:    &common.SDKTime{Time: time.Unix(200000, 0)},
					SizeInMBs:      types.Int64Ptr(2000),
				},
			},
		}, nil)

	computeService.EXPECT().
		DeleteImage(context.TODO(), core.DeleteImageRequest{ImageId: types.StringPtr("2")}).
		Return(core.DeleteImageResponse{}, nil)

	p.DeleteImage(ctx, "test")
}
