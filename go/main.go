package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func addSecretVersion(w io.Writer, parent string) error {
	// Declare the payload to store.
	payload := []byte("my super secret data")
	// Compute checksum, use Castagnoli polynomial. Providing a checksum
	// is optional.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(payload, crc32c))

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data:       payload,
			DataCrc32C: &checksum,
		},
	}

	// Call the API.
	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %w", err)
	}
	fmt.Fprintf(w, "Added secret version: %s\n", result.Name)
	return nil
}

func accessSecretVersion(w io.Writer, name string) (string, error) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "",fmt.Errorf("failed to access secret version: %w", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return "",fmt.Errorf("data corruption detected")
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	// fmt.Fprintln(w, string(result.Payload.Data))
	return string(result.Payload.Data), nil
}

type EnvData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	w := os.Stdout
	
	// parent := "projects/1061405048387/secrets/test-secret"
	// err := addSecretVersion(w, parent)
	// if err != nil {
	//     fmt.Fprintln(os.Stderr, "Unable to add secret version:", err)
	//     os.Exit(1)
	// }

	name := "projects/1061405048387/secrets/test-secret/versions/2"
	envData,err := accessSecretVersion(w, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to access secret version:", err)
		os.Exit(1)
	}

	var data EnvData
	err = json.Unmarshal([]byte(envData), &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to unmarshal json:", err)
		os.Exit(1)
	}

	fmt.Println(data.Username)
}
