package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func TestTickMultiCluster(t *testing.T) {
	t.Parallel()

	// For convenience - uncomment these as well as the "os" import
	// when doing local testing if you need to skip any sections.
	// os.Setenv("SKIP_", "true")
	// os.Setenv("TERRATEST_REGION", "us-east-1")
	// os.Setenv("SKIP_setup_ami", "true")
	// os.Setenv("SKIP_deploy_to_aws", "true")
	// os.Setenv("SKIP_validate_influxdb", "true")
	// os.Setenv("SKIP_validate_telegraf", "true")
	// os.Setenv("SKIP_validate_chronograf", "true")
	// os.Setenv("SKIP_validate_kapacitor", "true")
	// os.Setenv("SKIP_teardown", "true")

	var testcases = []struct {
		testName             string
		telegrafPackerInfo   PackerInfo
		influxdbPackerInfo   PackerInfo
		chronografPackerInfo PackerInfo
		kapacitorPackerInfo  PackerInfo
		sleepDuration        int
	}{
		{
			"TestTickMultiClusterUbuntu",
			PackerInfo{
				builderName:  "telegraf-ami-ubuntu",
				templatePath: "telegraf-ami/telegraf.json"},
			PackerInfo{
				builderName:  "influxdb-ami-ubuntu",
				templatePath: "influxdb-ami/influxdb.json"},
			PackerInfo{
				builderName:  "chronograf-ami-ubuntu",
				templatePath: "chronograf-ami/chronograf.json"},
			PackerInfo{
				builderName:  "kapacitor-ami-ubuntu",
				templatePath: "kapacitor-ami/kapacitor.json"},
			0,
		},
		{
			"TestTickMultiClusterAmazonLinux",
			PackerInfo{
				builderName:  "telegraf-ami-amazon-linux",
				templatePath: "telegraf-ami/telegraf.json"},
			PackerInfo{
				builderName:  "influxdb-ami-amazon-linux",
				templatePath: "influxdb-ami/influxdb.json"},
			PackerInfo{
				builderName:  "chronograf-ami-amazon-linux",
				templatePath: "chronograf-ami/chronograf.json"},
			PackerInfo{
				builderName:  "kapacitor-ami-amazon-linux",
				templatePath: "kapacitor-ami/kapacitor.json"},
			3,
		},
	}

	for _, testCase := range testcases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			// This is terrible - but attempt to stagger the test cases to
			// avoid a concurrency issue
			time.Sleep(time.Duration(testCase.sleepDuration) * time.Second)

			examplesDir := test_structure.CopyTerraformFolderToTemp(t, "..", "/examples")

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.Destroy(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "setup_ami", func() {
				awsRegion := aws.GetRandomRegion(t, nil, []string{"eu-north-1"})

				telegrafAmiID := buildAmi(t, fmt.Sprintf("%s/%s", examplesDir, testCase.telegrafPackerInfo.templatePath), testCase.telegrafPackerInfo.builderName, awsRegion)
				influxdbAmiID := buildAmi(t, fmt.Sprintf("%s/%s", examplesDir, testCase.influxdbPackerInfo.templatePath), testCase.influxdbPackerInfo.builderName, awsRegion)
				chronografAmiID := buildAmi(t, fmt.Sprintf("%s/%s", examplesDir, testCase.chronografPackerInfo.templatePath), testCase.chronografPackerInfo.builderName, awsRegion)
				kapacitorAmiID := buildAmi(t, fmt.Sprintf("%s/%s", examplesDir, testCase.kapacitorPackerInfo.templatePath), testCase.kapacitorPackerInfo.builderName, awsRegion)

				uniqueID := strings.ToLower(random.UniqueId())

				appServerName := fmt.Sprintf("tick-app-server-%s", uniqueID)
				metaClusterName := fmt.Sprintf("influxdb-meta-%s", uniqueID)
				dataClusterName := fmt.Sprintf("influxdb-data-%s", uniqueID)
				chronografServerName := fmt.Sprintf("chronograf-server-%s", uniqueID)
				kapacitorServerName := fmt.Sprintf("kapacitor-server-%s", uniqueID)

				keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
				test_structure.SaveEc2KeyPair(t, examplesDir, keyPair)

				licenseKey := os.Getenv("LICENSE_KEY")
				sharedSecret := os.Getenv("SHARED_SECRET")

				require.NotEmpty(t, licenseKey, "License key must be set as an env var and not included as plain-text")
				require.NotEmpty(t, sharedSecret, "Shared secret must be set as an env var and not included as plain-text")

				terraformOptions := &terraform.Options{
					// The path to where your Terraform code is located
					TerraformDir: fmt.Sprintf("%s/tick-multi-cluster", examplesDir),
					Vars: map[string]interface{}{
						"aws_region":                       awsRegion,
						"telegraf_ami_id":                  telegrafAmiID,
						"influxdb_ami_id":                  influxdbAmiID,
						"chronograf_ami_id":                chronografAmiID,
						"kapacitor_ami_id":                 kapacitorAmiID,
						"ssh_key_name":                     keyPair.Name,
						"app_server_name":                  appServerName,
						"influxdb_meta_nodes_cluster_name": metaClusterName,
						"influxdb_data_nodes_cluster_name": dataClusterName,
						"chronograf_server_name":           chronografServerName,
						"kapacitor_server_name":            kapacitorServerName,
						"license_key":                      licenseKey,
						"shared_secret":                    sharedSecret,
					},
				}

				test_structure.SaveTerraformOptions(t, examplesDir, terraformOptions)
			})

			test_structure.RunTestStage(t, "deploy_to_aws", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.InitAndApply(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "validate_influxdb", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "influxdb_dns")
				port := terraform.Output(t, terraformOptions, "influxdb_port")
				validateInfluxdb(t, endpoint, port)
			})

			test_structure.RunTestStage(t, "validate_telegraf", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "influxdb_dns")
				port := terraform.Output(t, terraformOptions, "influxdb_port")
				databaseName := terraform.Output(t, terraformOptions, "telegraf_database")
				validateTelegraf(t, endpoint, port, databaseName)
			})

			test_structure.RunTestStage(t, "validate_chronograf", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "chronograf_dns")
				port := terraform.Output(t, terraformOptions, "chronograf_port")
				validateChronograf(t, endpoint, port)
			})

			test_structure.RunTestStage(t, "validate_kapacitor", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "kapacitor_dns")
				port := terraform.Output(t, terraformOptions, "kapacitor_port")
				validateKapacitor(t, endpoint, port)
			})

		})
	}
}
