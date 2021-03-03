package oci

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"github.com/olekukonko/tablewriter"
	"github.com/oracle/oci-go-sdk/core"
)

// CreateInstance launch a server in oci using an existing image
func (p *Provider) CreateInstance(ctx *lepton.Context) error {
	image, err := p.getImageByName(ctx, ctx.Config().CloudConfig.ImageName)
	if err != nil {
		return err
	}

	_, err = p.computeClient.LaunchInstance(context.TODO(), core.LaunchInstanceRequest{
		LaunchInstanceDetails: core.LaunchInstanceDetails{
			AvailabilityDomain: types.StringPtr(""),
			CompartmentId:      types.StringPtr(""),
			Shape:              types.StringPtr(""),
			SourceDetails: core.InstanceSourceViaImageDetails{
				ImageId: &image.ID,
			},
		},
	})
	if err != nil {
		ctx.Logger().Error(err.Error())
		return errors.New("failed launching instance")
	}

	return nil
}

// ListInstances prints servers list managed by oci in table
func (p *Provider) ListInstances(ctx *lepton.Context) (err error) {
	instances, err := p.GetInstances(ctx)
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Status", "Private Ips", "Public Ips", "Image", "Created"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})

	table.SetRowLine(true)

	for _, i := range instances {
		var rows []string

		rows = append(rows, i.ID)
		rows = append(rows, i.Name)
		rows = append(rows, i.Status)
		rows = append(rows, strings.Join(i.PrivateIps, ", "))
		rows = append(rows, strings.Join(i.PublicIps, ", "))
		rows = append(rows, i.Image)
		rows = append(rows, i.Created)

		table.Append(rows)
	}

	table.Render()

	return
}

// GetInstances returns the list of servers managed by upcloud
func (p *Provider) GetInstances(ctx *lepton.Context) (instances []lepton.CloudInstance, err error) {
	instances = []lepton.CloudInstance{}

	result, err := p.computeClient.ListInstances(context.TODO(), core.ListInstancesRequest{})
	if err != nil {
		return
	}

	for _, i := range result.Items {
		instances = append(instances, lepton.CloudInstance{
			ID:      *i.CompartmentId,
			Name:    *i.DisplayName,
			Status:  string(i.LifecycleState),
			Created: lepton.Time2Human(i.TimeCreated.Time),
		})
	}

	return
}

// GetInstanceByID return a oci instance details with ID specified
func (p *Provider) GetInstanceByID(ctx *lepton.Context, id string) (instance *lepton.CloudInstance, err error) {
	return
}

// DeleteInstance removes an instance
func (p *Provider) DeleteInstance(ctx *lepton.Context, instancename string) error {
	return nil
}

// StopInstance stops an instance
func (p *Provider) StopInstance(ctx *lepton.Context, instancename string) error {
	return nil
}

// StartInstance starts an instance
func (p *Provider) StartInstance(ctx *lepton.Context, instancename string) error {
	return nil
}

// GetInstanceLogs returns instance log
func (p *Provider) GetInstanceLogs(ctx *lepton.Context, instancename string) (string, error) {
	return "", nil
}

// PrintInstanceLogs prints instances logs on console
func (p *Provider) PrintInstanceLogs(ctx *lepton.Context, instancename string, watch bool) error {
	return nil
}
