package test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"path/filepath"

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
		// {
		// 	"TestTickMultiClusterAmazonLinux",
		// 	PackerInfo{
		// 		builderName:  "telegraf-ami-amazon-linux",
		// 		templatePath: "telegraf-ami/telegraf.json"},
		// 	PackerInfo{
		// 		builderName:  "influxdb-ami-amazon-linux",
		// 		templatePath: "influxdb-ami/influxdb.json"},
		// 	PackerInfo{
		// 		builderName:  "chronograf-ami-amazon-linux",
		// 		templatePath: "chronograf-ami/chronograf.json"},
		// 	PackerInfo{
		// 		builderName:  "kapacitor-ami-amazon-linux",
		// 		templatePath: "kapacitor-ami/kapacitor.json"},
		// 	3,
		// },
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

			test_structure.RunTestStage(t, "setup_ami", func() {
				awsRegion := aws.GetRandomRegion(t, nil, []string{"eu-north-1"})

				influxAmis := buildAllAmis(
					t,
					awsRegion,
					&testCase.telegrafPackerInfo,
					&testCase.influxdbPackerInfo,
					&testCase.chronografPackerInfo,
					&testCase.kapacitorPackerInfo,
					examplesDir,
				)

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
						"telegraf_ami_id":                  influxAmis.TelegrafAmiID,
						"influxdb_ami_id":                  influxAmis.InfluxdbAmiID,
						"chronograf_ami_id":                influxAmis.ChronografAmiID,
						"kapacitor_ami_id":                 influxAmis.KapacitorAmiID,
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

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, examplesDir)
				terraform.Destroy(t, terraformOptions)

				// awsRegion := terraformOptions.Vars["aws_region"].(string)
				// keyPair := test_structure.LoadEc2KeyPair(t, examplesDir)

				// aws.DeleteAmi(t, awsRegion, terraformOptions.Vars["telegraf_ami_id"].(string))
				// aws.DeleteAmi(t, awsRegion, terraformOptions.Vars["influxdb_ami_id"].(string))
				// aws.DeleteAmi(t, awsRegion, terraformOptions.Vars["chronograf_ami_id"].(string))
				// aws.DeleteAmi(t, awsRegion, terraformOptions.Vars["kapacitor_ami_id"].(string))
				// aws.DeleteEC2KeyPair(t, keyPair)
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

func buildAllAmis(t *testing.T, awsRegion string, telegrafPackerInfo *PackerInfo, influxdbPackerInfo *PackerInfo, chronografPackerInfo *PackerInfo, kapacitorPackerInfo *PackerInfo, examplesDir string) *InfluxAmis {
	var waitForAmis sync.WaitGroup
	waitForAmis.Add(4)

	var telegrafAmiID string
	var influxdBAmiID string
	var chronografAmiID string
	var kapacitorAmiID string

	go func() {
		defer waitForAmis.Done()
		telegrafAmiID = buildAmi(t, filepath.Join(examplesDir, telegrafPackerInfo.templatePath), telegrafPackerInfo.builderName, awsRegion)
	}()
	go func() {
		defer waitForAmis.Done()
		influxdBAmiID = buildAmi(t, filepath.Join(examplesDir, influxdbPackerInfo.templatePath), influxdbPackerInfo.builderName, awsRegion)
	}()
	go func() {
		defer waitForAmis.Done()
		chronografAmiID = buildAmi(t, filepath.Join(examplesDir, chronografPackerInfo.templatePath), chronografPackerInfo.builderName, awsRegion)
	}()
	go func() {
		defer waitForAmis.Done()
		kapacitorAmiID = buildAmi(t, filepath.Join(examplesDir, kapacitorPackerInfo.templatePath), kapacitorPackerInfo.builderName, awsRegion)
	}()

	waitForAmis.Wait()

	if (telegrafAmiID == "") || (influxdBAmiID == "") || (chronografAmiID == "") || (kapacitorAmiID == "") {
		t.Fatalf("One of the AMIs was blank: telegraf:%s, influxdB:%s, chronograf:%s, kapacitor:%s", telegrafAmiID, influxdBAmiID, chronografAmiID, kapacitorAmiID)
	}

	return &InfluxAmis{
		TelegrafAmiID:   telegrafAmiID,
		InfluxdbAmiID:   influxdBAmiID,
		ChronografAmiID: chronografAmiID,
		KapacitorAmiID:  kapacitorAmiID,
	}
}

type InfluxAmis struct {
	TelegrafAmiID   string
	InfluxdbAmiID   string
	ChronografAmiID string
	KapacitorAmiID  string
}
