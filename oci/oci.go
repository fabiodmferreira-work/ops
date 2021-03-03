package oci

import (
	"context"

	"github.com/nanovms/ops/types"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"github.com/spf13/afero"
)

// ComputeService has OCI client methods to manage images and instances listing
type ComputeService interface {
	CreateImage(ctx context.Context, request core.CreateImageRequest) (response core.CreateImageResponse, err error)
	ListImages(ctx context.Context, request core.ListImagesRequest) (response core.ListImagesResponse, err error)
	DeleteImage(ctx context.Context, request core.DeleteImageRequest) (response core.DeleteImageResponse, err error)
	ListInstances(ctx context.Context, request core.ListInstancesRequest) (response core.ListInstancesResponse, err error)
	LaunchInstance(ctx context.Context, request core.LaunchInstanceRequest) (response core.LaunchInstanceResponse, err error)
}

// ComputeManagementService has OCI client methods to manage instances
type ComputeManagementService interface {
	CreateInstanceConfiguration(ctx context.Context, request core.CreateInstanceConfigurationRequest) (response core.CreateInstanceConfigurationResponse, err error)
	DeleteInstanceConfiguration(ctx context.Context, request core.DeleteInstanceConfigurationRequest) (response core.DeleteInstanceConfigurationResponse, err error)
	CreateInstancePool(ctx context.Context, request core.CreateInstancePoolRequest) (response core.CreateInstancePoolResponse, err error)
	StartInstancePool(ctx context.Context, request core.StartInstancePoolRequest) (response core.StartInstancePoolResponse, err error)
	StopInstancePool(ctx context.Context, request core.StopInstancePoolRequest) (response core.StopInstancePoolResponse, err error)
}

// StorageService has OCI client methods to manage storage block, required to upload images
type StorageService interface {
	PutObject(ctx context.Context, request objectstorage.PutObjectRequest) (response objectstorage.PutObjectResponse, err error)
}

// Provider has methods to interact with oracle cloud infrastructure
type Provider struct {
	computeClient           ComputeService
	computeManagementClient ComputeManagementService
	storageClient           StorageService
	fileSystem              afero.Fs
}

// NewProvider returns an instance of OCI Provider
func NewProvider() *Provider {
	return &Provider{
		computeClient: nil,
		storageClient: nil,
		fileSystem:    afero.NewOsFs(),
	}
}

// NewProviderWithClients returns an instance of OCI Provider with required clients initialized
func NewProviderWithClients(c ComputeService, cm ComputeManagementService, s StorageService, f afero.Fs) *Provider {
	return &Provider{c, cm, s, f}
}

// Initialize checks conditions to use oci
func (p *Provider) Initialize(providerConfig *types.ProviderConfig) (err error) {
	config := common.DefaultConfigProvider()

	p.computeClient, err = core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		return
	}

	p.storageClient, err = objectstorage.NewObjectStorageClientWithConfigurationProvider(config)
	if err != nil {
		return
	}

	p.computeManagementClient, err = core.NewComputeManagementClientWithConfigurationProvider(config)
	if err != nil {
		return
	}

	return
}
