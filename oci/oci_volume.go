package oci

import "github.com/nanovms/ops/lepton"

// CreateVolume creates a local volume and uploads the volume to oci
func (p *Provider) CreateVolume(ctx *lepton.Context, name, data, size, provider string) (vol lepton.NanosVolume, err error) {
	return
}

// GetAllVolumes returns a list of oci volumes
func (p *Provider) GetAllVolumes(ctx *lepton.Context) (vols *[]lepton.NanosVolume, err error) {
	return
}

// DeleteVolume removes an oci volume
func (p *Provider) DeleteVolume(ctx *lepton.Context, name string) error {
	return nil
}

// AttachVolume attaches a volume to an oci instance
func (p *Provider) AttachVolume(ctx *lepton.Context, image, name, mount string) error {
	return nil
}

// DetachVolume detaches a volume from an oci instance
func (p *Provider) DetachVolume(ctx *lepton.Context, image, name string) error {
	return nil
}
