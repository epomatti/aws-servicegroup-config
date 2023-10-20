package main

import (
	"context"
	"fmt"
	"log"
	"main/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	utils.Check(err)

	groupName := getGroupFromSSM(cfg)

	fmt.Println(groupName)

	// client := ec2.NewFromConfig(cfg)

	// &ec2.group
	// input := &ec2.StartInstancesInput{
	// 	InstanceIds: body.InstanceIds,
	// }

}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

func getGroupFromSSM(cfg aws.Config) string {
	var name = "/terraform/security-group-id"
	client := ssm.NewFromConfig(cfg)
	input := &ssm.GetParameterInput{
		Name: &name,
	}
	param, err := client.GetParameter(context.TODO(), input)
	utils.Check(err)

	return *param.Parameter.Value
}
