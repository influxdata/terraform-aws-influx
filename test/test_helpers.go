package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/gruntwork-io/terratest/modules/packer"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Addr:    fmt.Sprintf("http://%s:%s", endpoint, port),
		Timeout: time.Second * 60,
	})

	require.NoError(t, err, "Unable to connect to InfluxDB endpoint")

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
			logger.Logf(t, "Query failed: %s", response.Error().Error())
			return "", response.Error()
		}

		return "", nil
	})

	// Write to database
	branchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  databaseName,
		Precision: "s",
	})

	require.NoError(t, err, "Unable to create branch points")

	point, err := client.NewPoint(
		metric,
		map[string]string{"city": city},
		map[string]interface{}{"value": value},
		timestamp,
	)

	require.NoError(t, err, "Unable to create a point")

	branchPoints.AddPoint(point)
	err = c.Write(branchPoints)
	require.NoError(t, err, "Unable to write to database")

	// Read from database
	response, err := c.Query(client.Query{
		Command:  fmt.Sprintf("SELECT * FROM %s", metric),
		Database: databaseName,
	})

	require.NoError(t, err, "Unable to read from database")
	require.NoError(t, response.Error(), "Query failed")

	assert.Len(t, response.Results, 1)

	// Verify returned result
	series := response.Results[0].Series[0]

	assert.Equal(t, metric, series.Name)
	assert.Equal(t, city, series.Values[0][1])

	returnedValue, _ := series.Values[0][2].(json.Number).Int64()
	assert.Equal(t, value, returnedValue)
}
