package ansible

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/liangrog/taws/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultInventoryFileName = "ec2-inventory"
)

// get-inventory options
type FlagsGetInventory struct {
	// Filter by autoscaling group name string
	AsgNameFilter string

	// Output go file name
	InventoryFile string
}

// Check if k8s context is set
func (ops *FlagsGetInventory) HasAsgNameFilter() bool {
	return len(ops.AsgNameFilter) > 0
}

// Global
var OpsGetInventory *FlagsGetInventory

// Initiate get-inventory options
func NewOpsGetInventory() *FlagsGetInventory {
	// Current working directory
	cwd, err := os.Getwd()
	if err != nil {
		panic("Unable to get current working directory stats")
	}

	OpsGetInventory = new(FlagsGetInventory)
	OpsGetInventory.InventoryFile = filepath.Join(cwd, DefaultInventoryFileName)

	return OpsGetInventory
}

// Get inventory output as a file
func GetInventory() {
	// Get all instances under autoscaling

	// Place holder for instance ids
	var instIds []*string
	// Autoscaling group to instance ids mapping
	asInsts := make(map[string][]string)

	for {
		as := autoscaling.New(utils.AwsSess)
		asIn := &autoscaling.DescribeAutoScalingInstancesInput{}
		asOut, err := as.DescribeAutoScalingInstances(asIn)

		if err != nil {
			utils.ExitWithError(err)
		}

		for _, d := range asOut.AutoScalingInstances {
			// If set search context
			if OpsGetInventory.HasAsgNameFilter() &&
				!strings.Contains(*d.AutoScalingGroupName, OpsGetInventory.AsgNameFilter) {
				continue
			}

			asInsts[*d.AutoScalingGroupName] = append(asInsts[*d.AutoScalingGroupName], *d.InstanceId)
			instIds = append(instIds, d.InstanceId)
		}

		if asOut.NextToken == nil {
			break
		}
	}

	// Get instance private IPs
	idToIp := make(map[string]string)
	for {
		ec := ec2.New(utils.AwsSess)
		ecIn := &ec2.DescribeInstancesInput{
			InstanceIds: instIds,
		}
		ecOut, err := ec.DescribeInstances(ecIn)

		if err != nil {
			utils.ExitWithError(err)
		}

		for _, r := range ecOut.Reservations {
			for _, i := range r.Instances {
				idToIp[*i.InstanceId] = *i.PrivateIpAddress
			}
		}

		if ecOut.NextToken == nil {
			break
		}
	}

	// Inventory mapping: autoscaling group -> instance IPs
	var buf bytes.Buffer
	for k, v := range asInsts {
		buf.WriteString(fmt.Sprintf("[%s]\n", k))
		for _, v2 := range v {
			buf.WriteString(fmt.Sprintf("%s\n", idToIp[v2]))
		}

		buf.WriteString(fmt.Sprintf("\n"))
	}

	// Write to file
	err := ioutil.WriteFile(OpsGetInventory.InventoryFile, []byte(buf.String()), 0664)
	if err != nil {
		utils.ExitWithError(err)
	}
}
