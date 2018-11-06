package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestInfluxDBMultiCluster(t *testing.T) {
	t.Parallel()

	// For convenience - uncomment these as well as the "os" import
	// when doing local testing if you need to skip any sections.
	// os.Setenv("SKIP_", "true")
	// os.Setenv("TERRATEST_REGION", "us-east-1")
	// os.Setenv("SKIP_setup_ami", "true")
	// os.Setenv("SKIP_deploy_to_aws", "true")
	// os.Setenv("SKIP_validate", "true")
	// os.Setenv("SKIP_teardown", "true")

	examplesDir := test_structure.CopyTerraformFolderToTemp(t, "..", "/examples")
	amiDir := fmt.Sprintf("%s/influxdb-ami", examplesDir)

	var testcases = []struct {
		testName   string
		packerInfo PackerInfo
	}{
		{
			"TestInfluxDBMultiClusterUbuntu",
			PackerInfo{
				builderName:  "influxdb-ami-ubuntu",
				templatePath: fmt.Sprintf("%s/influxdb.json", amiDir)},
		},
		{
			"TestInfluxDBMultiClusterAmazonLinux",
			PackerInfo{
				builderName:  "influxdb-ami-amazon-linux",
				templatePath: fmt.Sprintf("%s/influxdb.json", amiDir)},
		},
	}

	for _, testCase := range testcases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			awsRegion := aws.GetRandomRegion(t, nil, nil)

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.Destroy(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "setup_ami", func() {
				amiID := buildAmi(t, testCase.packerInfo.templatePath, testCase.packerInfo.builderName, awsRegion)

				uniqueID := strings.ToLower(random.UniqueId())
				metaClusterName := fmt.Sprintf("influxdb-meta-%s", uniqueID)
				dataClusterName := fmt.Sprintf("influxdb-data-%s", uniqueID)

				keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
				test_structure.SaveEc2KeyPair(t, examplesDir, keyPair)

				terraformOptions := &terraform.Options{
					// The path to where your Terraform code is located
					TerraformDir: fmt.Sprintf("%s/influxdb-multi-cluster", examplesDir),
					Vars: map[string]interface{}{
						"aws_region":                       awsRegion,
						"ami_id":                           amiID,
						"ssh_key_name":                     keyPair.Name,
						"influxdb_meta_nodes_cluster_name": metaClusterName,
						"influxdb_data_nodes_cluster_name": dataClusterName,
						"license_key":                      os.Getenv("LICENSE_KEY"),
						"shared_secret":                    os.Getenv("SHARED_SECRET"),
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
