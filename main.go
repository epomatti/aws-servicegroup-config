package main

import (
	"context"
	"log"
	"main/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	utils.Check(err)

	groupName := getGroupFromSSM(cfg)
	revokeRules(cfg, groupName)
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

func revokeRules(cfg aws.Config, groupId string) {
	client := ec2.NewFromConfig(cfg)

	filterName := "group-id"
	filters := []types.Filter{{Name: &filterName, Values: []string{groupId}}}

	input := &ec2.DescribeSecurityGroupRulesInput{
		Filters: filters,
	}

	output, err := client.DescribeSecurityGroupRules(context.TODO(), input)
	utils.Check(err)

	rules := output.SecurityGroupRules
	ingress := make([]string, 0)

	for _, rule := range rules {
		if *rule.IsEgress == false {
			ingress = append(ingress, *rule.SecurityGroupRuleId)
		}
	}

	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		SecurityGroupRuleIds: ingress,
		GroupId:              &groupId,
	}

	_, err = client.RevokeSecurityGroupIngress(context.TODO(), revokeInput)
	utils.Check(err)
}

// func updateGroup(cfg aws.Config, groupName string) {
// 	client := ec2.NewFromConfig(cfg)

// 	input := &ec2.DescribeSecurityGroupRulesInput{
// 		groupName: groupName,
// 	}

// 	input := &ec2.StartInstancesInput{
// 		InstanceIds: body.InstanceIds,
// 	}
// }
