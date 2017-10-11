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

	FlagValueFilterByTags    = "tags"
	FlagValueFilterByAsgName = "asg-name"
	FlagValueGroupByASG      = "asg"
)

// get-inventory options
type FlagsGetInventory struct {
	// Filter
	// Available:
	// tags
	// asg-name
	FilterBy string

	// If filter by tags
	// what key and values
	// format key=value;key=value
	FilterValue string

	// Group by
	// Available: asg
	GroupBy string

	// Use Public IP
	UsePublicIp bool

	// Output go file name
	ToFile string
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
	OpsGetInventory.ToFile = filepath.Join(cwd, DefaultInventoryFileName)

	return OpsGetInventory
}

// Get inventory output as a file
func GetInventory() {
	// Get all instances under autoscaling
	ec2Input := new(ec2.DescribeInstancesInput)

	asgInsts := make(map[string][]string)
	if OpsGetInventory.GroupBy == FlagValueGroupByASG {
		asgInsts = GetInstanceIdsViaASG()
		var instIds []*string

		// Get all instance Ids into a list
		for _, g := range asgInsts {
			for idx, _ := range g {
				instIds = append(instIds, &g[idx])
			}
		}

		ec2Input.InstanceIds = instIds
	}

	insts := GetInstances(ec2Input)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n"))

	switch OpsGetInventory.GroupBy {
	case FlagValueGroupByASG:
		// Inventory mapping: autoscaling group -> instance IPs
		for k, v := range asgInsts {
			buf.WriteString(fmt.Sprintf("[%s]\n", k))
			for _, v2 := range v {
				buf.WriteString(fmt.Sprintf("%s\n", insts[v2]))
			}

			buf.WriteString(fmt.Sprintf("\n"))
		}
	default:
		for _, ip := range insts {
			buf.WriteString(fmt.Sprintf("%s\n", ip))
		}
	}

	// Write to file
	err := ioutil.WriteFile(OpsGetInventory.ToFile, []byte(buf.String()), 0664)
	if err != nil {
		utils.ExitWithError(err)
	}
}

// Get instance IDs via autoscaling group
func GetInstanceIdsViaASG() map[string][]string {
	// Autoscaling group to instance ids mapping
	insts := make(map[string][]string)

	for {
		as := autoscaling.New(utils.AwsSess)
		asIn := &autoscaling.DescribeAutoScalingInstancesInput{}
		asOut, err := as.DescribeAutoScalingInstances(asIn)

		if err != nil {
			utils.ExitWithError(err)
		}

		for _, d := range asOut.AutoScalingInstances {
			// If set search context
			if OpsGetInventory.FilterBy == FlagValueFilterByAsgName &&
				!strings.Contains(*d.AutoScalingGroupName, OpsGetInventory.FilterValue) {
				continue
			}

			insts[*d.AutoScalingGroupName] = append(insts[*d.AutoScalingGroupName], *d.InstanceId)
		}

		if asOut.NextToken == nil {
			break
		}
	}

	return insts
}

// Get EC2 instances
func GetInstances(input *ec2.DescribeInstancesInput) map[string]string {
	// Get instance IPs
	idToIp := make(map[string]string)

	for {
		ec := ec2.New(utils.AwsSess)
		ecOut, err := ec.DescribeInstances(input)

		if err != nil {
			utils.ExitWithError(err)
		}

		for _, r := range ecOut.Reservations {
			for _, i := range r.Instances {
				// If there is tag filter
				if OpsGetInventory.FilterBy == FlagValueFilterByTags &&
					!FilterTags(i.Tags, OpsGetInventory.FilterValue) {
					continue
				}

				// If use public IP
				if OpsGetInventory.UsePublicIp {
					if i.PublicIpAddress != nil {
						idToIp[*i.InstanceId] = *i.PublicIpAddress
					} else {
						fmt.Printf("Warning: instance %s, private IP %s has no public IP address, ignored from inventory\n", *i.InstanceId, *i.PrivateIpAddress)
						continue
					}
				} else {
					idToIp[*i.InstanceId] = *i.PrivateIpAddress
				}
			}
		}

		if ecOut.NextToken == nil {
			break
		}
	}

	return idToIp
}

// Check if tags contains all tag filter value
func FilterTags(tags []*ec2.Tag, fv string) bool {
	tagSets := strings.Split(fv, ";")
	total := len(tagSets)
	count := 0
	for _, ts := range tagSets {
		tv := strings.Split(ts, "=")
		for _, t := range tags {
			if *t.Key == tv[0] &&
				*t.Value == tv[1] {
				count++
			}
		}
	}

	if total == count {
		return true
	}

	return false
}
