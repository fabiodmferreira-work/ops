package oci_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nanovms/ops/mock_oci"
	"github.com/nanovms/ops/oci"
	"github.com/spf13/afero"
)

func NewProvider(t *testing.T) (*oci.Provider, *mock_oci.MockComputeService, *mock_oci.MockComputeManagementService, *mock_oci.MockStorageService, afero.Fs) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	computeClient := mock_oci.NewMockComputeService(ctrl)
	computeManagementClient := mock_oci.NewMockComputeManagementService(ctrl)
	storageClient := mock_oci.NewMockStorageService(ctrl)
	fileReader := afero.NewMemMapFs()

	return oci.NewProviderWithClients(computeClient, computeManagementClient, storageClient, fileReader), computeClient, computeManagementClient, storageClient, fileReader
}
