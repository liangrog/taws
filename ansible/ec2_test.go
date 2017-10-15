package ansible

import (
	"bufio"
	"encoding/json"
	_ "fmt"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	asgTestData   = "../testdata/asg.json"
	ec2TestData   = "../testdata/ec2.json"
	inventoryFile = "../testdata/inventory"
)

// Mocking ASG service
type msvcAsg struct {
	autoscalingiface.AutoScalingAPI
}

func (*msvcAsg) DescribeAutoScalingInstances(*autoscaling.DescribeAutoScalingInstancesInput) (*autoscaling.DescribeAutoScalingInstancesOutput, error) {
	raw, _ := ioutil.ReadFile(asgTestData)

	var insts []*autoscaling.InstanceDetails
	json.Unmarshal([]byte(raw), &insts)

	asgOut := new(autoscaling.DescribeAutoScalingInstancesOutput)
	asgOut.AutoScalingInstances = insts

	return asgOut, nil
}

func TestGetInstanceIdsViaASG(t *testing.T) {
	svcAsg = new(msvcAsg)
	OpsGetInventory = new(FlagsGetInventory)

	// No filter
	result := GetInstanceIdsViaASG()
	assert.Equal(t, len(result), 2)

	// With filter
	OpsGetInventory.FilterBy = FlagValueFilterByAsgName
	OpsGetInventory.FilterValue = "node"
	result = GetInstanceIdsViaASG()
	assert.Equal(t, len(result), 1)
}

// Mocking EC2 service
type msvcEc2 struct {
	ec2iface.EC2API
}

func (*msvcEc2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	raw, _ := ioutil.ReadFile(ec2TestData)

	var res []*ec2.Reservation
	json.Unmarshal([]byte(raw), &res)

	out := new(ec2.DescribeInstancesOutput)
	out.Reservations = res

	return out, nil
}

func TestGetInstances(t *testing.T) {
	svcEc2 = new(msvcEc2)
	OpsGetInventory = new(FlagsGetInventory)

	// Without filter
	result := GetInstances(new(ec2.DescribeInstancesInput))
	assert.Equal(t, len(result), 2)

	// With filter
	OpsGetInventory.FilterBy = FlagValueFilterByTags
	OpsGetInventory.FilterValue = "Name=nodes.k8s"
	result = GetInstances(new(ec2.DescribeInstancesInput))
	assert.Equal(t, len(result), 1)

	// Use public IP
	OpsGetInventory.FilterBy = ""
	OpsGetInventory.FilterValue = ""
	OpsGetInventory.UsePublicIp = true
	result = GetInstances(new(ec2.DescribeInstancesInput))
	assert.Equal(t, len(result), 1)
}

func TestFilterTags(t *testing.T) {
	k1 := "key1"
	v1 := "value1"
	k2 := "key2"
	v2 := "value2"

	tags := []*ec2.Tag{
		{
			Key:   &k1,
			Value: &v1,
		},
		{
			Key:   &k2,
			Value: &v2,
		},
	}

	assert.True(t, FilterTags(tags, "key2=value2"))
	assert.True(t, FilterTags(tags, "key1=value1;key2=value2"))
	assert.False(t, FilterTags(tags, "key3=value3"))
}

// Simply testing the output file as all scenarios
// have already been testing in other unit tests
func TestGetInventory(t *testing.T) {
	svcAsg = new(msvcAsg)
	svcEc2 = new(msvcEc2)

	OpsGetInventory = new(FlagsGetInventory)
	OpsGetInventory.ToFile = inventoryFile

	GetInventory()

	_, err := os.Stat(inventoryFile)
	assert.Equal(t, err, nil)

	// Count file lines, should be two IP addresses
	lc := 0
	f, err := os.Open(inventoryFile)
	defer f.Close()
	defer os.Remove(inventoryFile)
	s := bufio.NewScanner(f)
	for s.Scan() {
		if len(s.Text()) > 0 {
			lc++
		}
	}
	assert.Equal(t, lc, 2)
}
