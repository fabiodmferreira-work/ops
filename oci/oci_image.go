package oci

import (
	"context"
	"errors"
	"fmt"

	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// BuildImage creates local image
func (p *Provider) BuildImage(ctx *lepton.Context) (string, error) {
	c := ctx.Config()
	err := lepton.BuildImage(*c)
	if err != nil {
		return "", err
	}

	return c.RunConfig.Imagename, nil
}

// BuildImageWithPackage creates local image using package image
func (p *Provider) BuildImageWithPackage(ctx *lepton.Context, pkgpath string) (string, error) {
	c := ctx.Config()
	err := lepton.BuildImageFromPackage(pkgpath, *c)
	if err != nil {
		return "", err
	}

	return ctx.Config().RunConfig.Imagename, nil
}

// CreateImage creates a storage object and upload image
func (p *Provider) CreateImage(ctx *lepton.Context, imagePath string) (err error) {
	bucketNamespace := ctx.Config().CloudConfig.BucketNamespace
	bucketName := ctx.Config().CloudConfig.BucketName
	imageName := ctx.Config().CloudConfig.ImageName

	image, err := p.fileSystem.Open(imagePath)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return fmt.Errorf("failed reading file %s", imagePath)
	}

	imageStats, err := image.Stat()
	if err != nil {
		ctx.Logger().Error(err.Error())
		return fmt.Errorf("failed getting file stats of %s", imagePath)
	}

	imageSize := imageStats.Size()

	_, err = p.storageClient.PutObject(context.TODO(), objectstorage.PutObjectRequest{
		NamespaceName: &bucketNamespace,
		BucketName:    &bucketName,
		ContentLength: &imageSize,
		ObjectName:    &imageName,
		PutObjectBody: image,
	})
	if err != nil {
		ctx.Logger().Error(err.Error())
		return errors.New("failed uploading image")
	}

	_, err = p.computeClient.CreateImage(context.TODO(), core.CreateImageRequest{
		CreateImageDetails: core.CreateImageDetails{
			ImageSourceDetails: core.ImageSourceViaObjectStorageTupleDetails{
				NamespaceName:   &bucketNamespace,
				BucketName:      &bucketName,
				ObjectName:      &imageName,
				SourceImageType: core.ImageSourceDetailsSourceImageTypeQcow2,
			},
		},
	})
	if err != nil {
		ctx.Logger().Error(err.Error())
		return errors.New("failed importing image from storage")
	}

	return
}

// ListImages prints oci images in table format
func (p *Provider) ListImages(ctx *lepton.Context) error {
	return nil
}

// GetImages returns the list of images
func (p *Provider) GetImages(ctx *lepton.Context) (images []lepton.CloudImage, err error) {
	images = []lepton.CloudImage{}

	imagesList, err := p.computeClient.ListImages(context.TODO(), core.ListImagesRequest{})
	if err != nil {
		ctx.Logger().Error(err.Error())
		return nil, errors.New("failed getting images")
	}

	for _, i := range imagesList.Items {
		images = append(images, lepton.CloudImage{
			ID:      *i.CompartmentId,
			Name:    *i.DisplayName,
			Status:  string(i.LifecycleState),
			Created: i.TimeCreated.Time,
			Size:    *i.SizeInMBs / 1024,
		})
	}

	return
}

// DeleteImage removes oci image
func (p *Provider) DeleteImage(ctx *lepton.Context, imagename string) (err error) {

	image, err := p.getImageByName(ctx, imagename)
	if err != nil {
		return
	}

	_, err = p.computeClient.DeleteImage(context.TODO(), core.DeleteImageRequest{ImageId: &image.ID})

	return
}

func (p *Provider) getImageByName(ctx *lepton.Context, name string) (*lepton.CloudImage, error) {
	images, err := p.GetImages(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range images {
		if i.Name == name {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("image with name %s not found", name)
}

// ResizeImage is a stub
func (p *Provider) ResizeImage(ctx *lepton.Context, imagename string, hbytes string) error {
	return errors.New("Unsupported")
}

// SyncImage is a stub
func (p *Provider) SyncImage(config *types.Config, target lepton.Provider, imagename string) error {
	return errors.New("Unsupported")
}

// CustomizeImage is a stub
func (p *Provider) CustomizeImage(ctx *lepton.Context) (string, error) {
	return "", errors.New("Unsupported")
}
