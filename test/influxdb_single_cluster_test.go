package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"

	"github.com/influxdata/influxdb/client/v2"
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

	rootDir := test_structure.CopyTerraformFolderToTemp(t, "..", "/")
	examplesDir := fmt.Sprintf("%s/examples", rootDir)
	amiDir := fmt.Sprintf("%s/influxdb-ami", examplesDir)

	var testcases = []struct {
		testName   string
		packerInfo PackerInfo
	}{
		{
			"TestInfluxDBSingleClusterUbuntu",
			PackerInfo{
				builderName:  "influxdb-ami-ubuntu",
				templatePath: fmt.Sprintf("%s/influxdb.json", amiDir)},
		},
		{
			"TestInfluxDBSingleClusterAmazonLinux",
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
				clusterName := fmt.Sprintf("influxdb-%s", uniqueID)

				keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
				test_structure.SaveEc2KeyPair(t, examplesDir, keyPair)

				terraformOptions := &terraform.Options{
					// The path to where your Terraform code is located
					TerraformDir: fmt.Sprintf("%s", rootDir),
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

				databaseName := "automatedtest"
				metric := "temperature"
				city := "Aurora"
				value := int64(50)
				timestamp := time.Now()

				c, err := client.NewHTTPClient(client.HTTPConfig{
					Addr: fmt.Sprintf("http://%s:%s", endpoint, port),
				})

				if err != nil {
					t.Fatal("Unable to connect to InfluxDB endpoint")
				}

				defer c.Close()

				// Create database
				response, _ := c.Query(client.Query{
					Command: fmt.Sprintf("CREATE DATABASE %s", databaseName),
				})

				if response.Error() != nil {
					t.Fatalf("Query failed: %s", response.Error().Error())
				}

				// Write to database
				branchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
					Database:  databaseName,
					Precision: "s",
				})

				if err != nil {
					t.Fatal("Unable to create branch points")
				}

				point, err := client.NewPoint(
					metric,
					map[string]string{"city": city},
					map[string]interface{}{"value": value},
					timestamp,
				)

				if err != nil {
					t.Fatal("Unable to create a point")
				}

				branchPoints.AddPoint(point)
				err = c.Write(branchPoints)

				if err != nil {
					t.Fatal("Unable to write to database")
				}

				// Read from database
				response, _ = c.Query(client.Query{
					Command:  fmt.Sprintf("SELECT * FROM %s", metric),
					Database: databaseName,
				})

				if response.Error() != nil {
					t.Fatalf("Query failed: %s", response.Error().Error())
				}

				if len(response.Results) != 1 {
					t.Fatal("Was only expecting one result object")
				}

				// Verify returned result
				series := response.Results[0].Series[0]
				returnedValue, _ := series.Values[0][2].(json.Number).Int64()
				if series.Name != metric || series.Values[0][1] != city || returnedValue != value {
					t.Fatal("Invalid entry retrieved")
				}
			})
		})
	}
}
