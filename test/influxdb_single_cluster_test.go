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
	"github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestInfluxDBSingleCluster(t *testing.T) {
	t.Parallel()

	// For convenience - uncomment these as well as the "os" import
	// when doing local testing if you need to skip any sections.
	// os.Setenv("SKIP_", "true")
	// os.Setenv("TERRATEST_REGION", "us-east-1")
	// os.Setenv("SKIP_setup_ami", "true")
	// os.Setenv("SKIP_deploy_to_aws", "true")
	// os.Setenv("SKIP_validate", "true")
	// os.Setenv("SKIP_teardown", "true")

	var testcases = []struct {
		testName      string
		packerInfo    PackerInfo
		sleepDuration int
	}{
		{
			"TestInfluxDBSingleClusterUbuntu",
			PackerInfo{
				builderName:  "influxdb-ami-ubuntu",
				templatePath: "influxdb.json"},
			0,
		},
		{
			"TestInfluxDBSingleClusterAmazonLinux",
			PackerInfo{
				builderName:  "influxdb-ami-amazon-linux",
				templatePath: "influxdb.json"},
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

			rootDir := test_structure.CopyTerraformFolderToTemp(t, "..", "/")
			examplesDir := fmt.Sprintf("%s/examples", rootDir)
			amiDir := fmt.Sprintf("%s/influxdb-ami", examplesDir)
			templatePath := fmt.Sprintf("%s/%s", amiDir, testCase.packerInfo.templatePath)

			awsRegion := aws.GetRandomRegion(t, nil, nil)

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.Destroy(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "setup_ami", func() {
				amiID := buildAmi(t, templatePath, testCase.packerInfo.builderName, awsRegion)

				uniqueID := strings.ToLower(random.UniqueId())
				clusterName := fmt.Sprintf("influxdb-%s", uniqueID)

				keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
				test_structure.SaveEc2KeyPair(t, examplesDir, keyPair)

				terraformOptions := &terraform.Options{
					// The path to where your Terraform code is located
					TerraformDir: rootDir,
					Vars: map[string]interface{}{
						"aws_region":            awsRegion,
						"ami_id":                amiID,
						"ssh_key_name":          keyPair.Name,
						"influxdb_cluster_name": clusterName,
						"license_key":           os.Getenv("LICENSE_KEY"),
						"shared_secret":         os.Getenv("SHARED_SECRET"),
					},
				}

				test_structure.SaveTerraformOptions(t, examplesDir, terraformOptions)
			})

			test_structure.RunTestStage(t, "deploy_to_aws", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.InitAndApply(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "validate", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				endpoint := terraform.Output(t, terraformOptions, "lb_dns_name")
				port := terraform.Output(t, terraformOptions, "load_balancer_port")
				validateInfluxdb(t, endpoint, port)
			})
		})
	}
}
