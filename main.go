package main

import (
	"context"
	"fmt"
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

	groupId := getGroupFromSSM(cfg)
	revokeRules(cfg, groupId)
	authorizeRules(cfg, groupId)
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

	filterName := aws.String("group-id")
	filters := []types.Filter{{Name: filterName, Values: []string{groupId}}}

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

	if len(ingress) == 0 {
		return
	}

	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		SecurityGroupRuleIds: ingress,
		GroupId:              &groupId,
	}

	_, err = client.RevokeSecurityGroupIngress(context.TODO(), revokeInput)
	utils.Check(err)
}

func authorizeRules(cfg aws.Config, groupId string) {
	admins := utils.ReadYaml()
	client := ec2.NewFromConfig(cfg)
	port := aws.Int32(22)
	ipProtocol := aws.String("tcp")
	ipPermissions := make([]types.IpPermission, 0)
	ranges := make([]types.IpRange, 0)
	for _, admin := range admins {
		for _, ip := range admin.CidrBlocks {
			ipRange := types.IpRange{
				// FIXME: This was bugging
				// CidrIp: &ip,
				CidrIp:      aws.String(ip),
				Description: aws.String(admin.Name),
			}
			ranges = append(ranges, ipRange)
		}
	}
	permission := types.IpPermission{
		FromPort:   port,
		ToPort:     port,
		IpProtocol: ipProtocol,
		IpRanges:   ranges,
	}
	ipPermissions = append(ipPermissions, permission)
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       &groupId,
		IpPermissions: ipPermissions,
	}
	authorizeResponse, err := client.AuthorizeSecurityGroupIngress(context.TODO(), input)
	utils.Check(err)

	fmt.Println(authorizeResponse.Return)
}
