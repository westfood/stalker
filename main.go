package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {

	table()

}

func table() {

	data := [][]string{
		[]string{"Alfréd", "15", "10/20", "(10.32, 56.21, 30.25)"},
		[]string{"Belzezub", "30", "30/50", "(1,1,1)"},
		[]string{"Hortenz", "21", "80/80", "(1,1,1)"},
		[]string{"Pokey", "8", "30/40", "(1,1,1)"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NPC", "Rychlost", "Síla", "Místo"})
	table.AppendBulk(data)
	table.Render()
}

func printArgs() {

	args := os.Args[1:]
	fmt.Println(args)
}

func describeInstance() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ec2svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: []*string{aws.String("172.31.3.8")},
			},
		},
	}
	result, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println(result)
	}
	fmt.Println(params)

}
