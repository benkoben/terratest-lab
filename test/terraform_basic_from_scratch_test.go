package test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/stretchr/testify/assert"
)

// To keep it simple just use global variables to control
// the test.
var (
	storageAccountName                 string
	storageAccountResourceGroup        string
	subscriptionId                     string
	backendStorageAccountName          string
	backendStorageAccountResourceGroup string
)

// read input from environment
func getEnv() error {
	if _, ok := os.LookupEnv("TEST_STORAGE_ACCOUNT_NAME"); ok {
		storageAccountName = os.Getenv("TEST_STORAGE_ACCOUNT_NAME")
	} else {
		return errors.New("could not find environment variable TEST_STORAGE_ACCOUNT_NAME")
	}

	if _, ok := os.LookupEnv("TEST_RESOURCE_GROUP"); ok {
		storageAccountResourceGroup = os.Getenv("TEST_RESOURCE_GROUP")
	} else {
		return errors.New("could not find environment variable TEST_RESOURCE_GROUP")
	}

	if _, ok := os.LookupEnv("TEST_SUBSCRIPTION_ID"); ok {
		subscriptionId = os.Getenv("TEST_SUBSCRIPTION_ID")
	} else {
		return errors.New("could not find environment variable TEST_SUBSCRIPTION_ID")
	}

	if _, ok := os.LookupEnv("TEST_BACKEND_STORAGE_ACCOUNT"); ok {
		backendStorageAccountName = os.Getenv("TEST_BACKEND_STORAGE_ACCOUNT")
	} else {
		return errors.New("could not find environment variable TEST_BACKEND_STORAGE_ACCOUNT")
	}

	if _, ok := os.LookupEnv("TEST_BACKEND_RESOURCE_GROUP"); ok {
		backendStorageAccountResourceGroup = os.Getenv("TEST_BACKEND_RESOURCE_GROUP")
	} else {
		return errors.New("could not find environment variable TEST_BACKEND_RESOURCE_GROUP")
	}
	return nil
}

func TestTerraformBasicFromScratch(t *testing.T) {
	// read environment variables
	err := getEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/terraform-basic-from-scratch",
		Vars: map[string]interface{}{
			"resource_group_name":  storageAccountResourceGroup,
			"storage_account_name": storageAccountName,
		},
		BackendConfig: map[string]interface{}{
			"resource_group_name":  backendStorageAccountResourceGroup,
			"storage_account_name": backendStorageAccountName,
			"container_name":       "modstorage",
			"key":                  "modstorage.tfstate",
		},
		Reconfigure: true,
		Upgrade:     true,
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.Init(t, terraformOptions)
	// Test the module's idempotency
	terraform.ApplyAndIdempotent(t, terraformOptions)
	// Gather outputs from Apply stage
	output_storage_account_name := terraform.Output(t, terraformOptions, "storage_account_name")
	output_container_name := terraform.Output(t, terraformOptions, "container_name")
	// Assert if the outputs are as expected
	assert.Equal(t, storageAccountName, output_storage_account_name)
	assert.Equal(t, azure.StorageBlobContainerExists(
		t,
		output_container_name,
		output_storage_account_name,
		storageAccountResourceGroup,
		subscriptionId), true)
}
