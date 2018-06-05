package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net"
	"os"
	"strings"
)

func main() {

	handleCommand()

}

func handleCommand() {

	args := (os.Args[1:])
	size := len(args)

	switch size {

	case 0:

		fmt.Println("Nothing is sub optimal amount of arguments, gimme some please.")
		printHelp()

	case 1:

		addr := net.ParseIP(strings.Join(args, ""))

		if addr == nil {

			hostname := strings.Join(args, "")

			if strings.Contains(hostname, "compute.internal") {
				describeInstance(hostname, "private-dns-name")
			} else if strings.Contains(hostname, "i-") {
				describeInstance(hostname, "instance-id")
			}
		} else {
			describeInstance(addr.String(), "private-ip-address")
		}

	default:

		fmt.Println("Command supports only one argument.")
		printHelp()

	}
}

func printHelp() {
	fmt.Println("Return public hostname of AWS instance from instance id, private IP or hostname. Example: stalker 172.1.24.45")
}

func describeInstance(ip, filter string) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ec2svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(filter),
				Values: []*string{aws.String(ip)},
			},
		},
	}

	result, err := ec2svc.DescribeInstances(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, network := range instance.NetworkInterfaces {
				fmt.Println(*network.Association.PublicDnsName)
			}
		}
	}

	// if err != nil {
	// 	if awsErr, ok := err.(awserr.Error); ok {
	// 		fmt.Println("Error", awsErr)
	// 	} else {
	// 		fmt.Println("Error", err)
	// 	}
	// } else {
	// 	for _, reservation := range result.Reservations {
	// 		for _, instance := range reservation.Instances {
	// 			for _, network := range instance.NetworkInterfaces {
	// 				fmt.Println(*network.Association.PublicDnsName)
	// 			}
	// 		}
	// 	}
	// }
}
