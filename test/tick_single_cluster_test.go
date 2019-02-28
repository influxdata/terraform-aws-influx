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

func TestTickSingleCluster(t *testing.T) {
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
	// os.Setenv("SKIP_teardown", "true")

	var testcases = []struct {
		testName      string
		packerInfo    PackerInfo
		sleepDuration int
	}{
		{
			"TestTickSingleClusterUbuntu",
			PackerInfo{
				builderName:  "tick-ami-ubuntu",
				templatePath: "tick.json"},
			0,
		},
		{
			"TestTickSingleClusterAmazonLinux",
			PackerInfo{
				builderName:  "tick-ami-amazon-linux",
				templatePath: "tick.json"},
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
			amiDir := fmt.Sprintf("%s/tick-ami", examplesDir)
			templatePath := fmt.Sprintf("%s/%s", amiDir, testCase.packerInfo.templatePath)

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.Destroy(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "setup_ami", func() {
				awsRegion := aws.GetRandomRegion(t, nil, []string{"eu-north-1"})
				amiID := buildAmi(t, templatePath, testCase.packerInfo.builderName, awsRegion)

				uniqueID := strings.ToLower(random.UniqueId())
				clusterName := fmt.Sprintf("tick-%s", uniqueID)

				keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
				test_structure.SaveEc2KeyPair(t, examplesDir, keyPair)

				licenseKey := os.Getenv("LICENSE_KEY")
				sharedSecret := os.Getenv("SHARED_SECRET")

				require.NotEmpty(t, licenseKey, "License key must be set as an env var and not included as plain-text")
				require.NotEmpty(t, sharedSecret, "Shared secret must be set as an env var and not included as plain-text")

				terraformOptions := &terraform.Options{
					// The path to where your Terraform code is located
					TerraformDir: fmt.Sprintf("%s/tick-single-cluster", examplesDir),
					Vars: map[string]interface{}{
						"aws_region":    awsRegion,
						"ami_id":        amiID,
						"ssh_key_name":  keyPair.Name,
						"cluster_name":  clusterName,
						"license_key":   licenseKey,
						"shared_secret": sharedSecret,
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
				endpoint := terraform.Output(t, terraformOptions, "lb_dns_name")
				port := terraform.Output(t, terraformOptions, "influxdb_port")
				validateInfluxdb(t, endpoint, port)
			})

			test_structure.RunTestStage(t, "validate_telegraf", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "lb_dns_name")
				port := terraform.Output(t, terraformOptions, "influxdb_port")
				validateTelegraf(t, endpoint, port, "telegraf")
			})

			test_structure.RunTestStage(t, "validate_chronograf", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "lb_dns_name")
				port := terraform.Output(t, terraformOptions, "chronograf_port")
				validateChronograf(t, endpoint, port)
			})
		})
	}
}
