package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroupNames struct {
	Names []string
}

// Set appends the given value to the slice of names in the SecurityGroupNames struct.
//
// value: The value to be appended to the slice.
// error: If an error occurs while appending the value.
func (s *SecurityGroupNames) Set(value string) error {
	s.Names = append(s.Names, value)
	return nil
}

// String returns a string representation of the SecurityGroupNames struct.
//
// It joins the names of the security groups in the struct using a comma as the separator.
// The resulting string is returned.
func (s *SecurityGroupNames) String() string {
	return strings.Join(s.Names, ",")
}

// main is the entry point of the program.
//
// It creates a flag to specify the security group names.
// It parses the command line arguments.
// For each security group name, it gets the network interfaces that are attached to it.
// It prints the security group name and the network interfaces that are attached to it.
//
// No parameters.
// No return values.
func main() {
	// Create a flag to specify the security group names

	// Create a flag to specify the security group names
	var securityGroupNames SecurityGroupNames
	flag.Var(&securityGroupNames, "security-group-names", "The names of the security groups to include in the output")

	// Parse the command line arguments
	flag.Parse()

	// For each security group name, get the network interfaces that are attached to it
	for _, securityGroupName := range securityGroupNames.Names {
		networkInterfaces := getNetworkInterfacesForSecurityGroup(securityGroupName)
		// Print the security group name and the network interfaces that are attached to it
		fmt.Printf("Security group name: %s\n", securityGroupName)
		for _, networkInterface := range networkInterfaces {
			fmt.Printf("Network interfaces:\n")
			fmt.Printf("  NetworkInterface ID: %s\n", *networkInterface.NetworkInterfaceId)
			if networkInterface.Attachment != nil && networkInterface.Attachment.InstanceId != nil {
				fmt.Printf("  InstanceId: %s\n", *networkInterface.Attachment.InstanceId)
			}
			fmt.Printf("  Status: %s\n", networkInterface.Status)
			fmt.Println()
		}
	}
}

// getSecurityGroupNames retrieves the names of all security groups.
//
// It does so by creating a config using the LoadDefaultConfig function from the AWS SDK for Go.
// If an error occurs during the creation of the config, it panics.
//
// It then creates a context using the TODO function from the context package.
//
// Next, it creates an EC2 client using the NewFromConfig function from the AWS SDK for Go.
//
// After that, it describes the security groups using the DescribeSecurityGroupsInput struct from the AWS SDK for Go.
//
// If an error occurs during the execution of the DescribeSecurityGroups function, it panics.
//
// Finally, it retrieves the security group names by iterating over the security groups in the DescribeSecurityGroupsOutput struct and appending their names to a slice.
//
// The function returns a slice of strings containing the security group names.
func getSecurityGroupNames() []string {
	// Create a config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	// context
	ctx := context.TODO()

	// Create an EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Describe the security groups
	describeSecurityGroupsInput := &ec2.DescribeSecurityGroupsInput{}

	describeSecurityGroupsOutput, err := ec2Client.DescribeSecurityGroups(ctx, describeSecurityGroupsInput)

	if err != nil {
		panic(err)
	}

	// Get the security group names
	securityGroupNames := []string{}
	for _, securityGroup := range describeSecurityGroupsOutput.SecurityGroups {
		securityGroupNames = append(securityGroupNames, *securityGroup.GroupName)
	}

	return securityGroupNames
}

// getNetworkInterfacesForSecurityGroup retrieves the network interfaces for a given security group.
//
// securityGroupName: The name of the security group.
// []types.NetworkInterface: An array of network interfaces.
func getNetworkInterfacesForSecurityGroup(securityGroupName string) []types.NetworkInterface {
	// Create a config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// context
	ctx := context.TODO()

	// Create an EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Describe the network interfaces
	describeNetworkInterfacesOutput, err := ec2Client.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []string{securityGroupName},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Get the network interfaces
	networkInterfaces := []types.NetworkInterface{}
	for _, networkInterface := range describeNetworkInterfacesOutput.NetworkInterfaces {
		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces
}
