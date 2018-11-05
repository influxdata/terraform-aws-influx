package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer"
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
