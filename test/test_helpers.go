package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/gruntwork-io/terratest/modules/packer"
	client "github.com/influxdata/influxdb/client/v2"
)

type PackerInfo struct {
	templatePath string
	builderName  string
}

func buildAmi(t *testing.T, templatePath string, builderName string, awsRegion string) string {
	options := &packer.Options{
		Template: templatePath,
		Only:     builderName,
		Vars: map[string]string{
			"aws_region": awsRegion,
		},
	}

	return packer.BuildAmi(t, options)
}

func validateInfluxdb(t *testing.T, endpoint string, port string) {
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

	maxRetries := 15
	sleepBetweenRetries := 5 * time.Second

	// Create database
	retry.DoWithRetry(t, "Querying database", maxRetries, sleepBetweenRetries, func() (string, error) {
		response, err := c.Query(client.Query{
			Command: fmt.Sprintf("CREATE DATABASE %s", databaseName),
		})

		if err != nil {
			t.Logf("Query failed: %s", err.Error())
			return "", err
		}

		if response.Error() != nil {
			t.Logf("Query failed: %s", response.Error().Error())
			return "", response.Error()
		}

		return "", nil
	})

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
	response, err := c.Query(client.Query{
		Command:  fmt.Sprintf("SELECT * FROM %s", metric),
		Database: databaseName,
	})

	if err != nil {
		t.Fatal("Unable to read from database")
	}

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
		t.Fatal("Incorrect entry retrieved")
	}
}
