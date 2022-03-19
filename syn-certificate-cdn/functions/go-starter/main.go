package main

import (
    "context"
    "fmt"
    "github.com/tencentyun/scf-go-lib/cloudfunction"
)

type DefineEvent struct {
    SecretId string `json:"secret_id"`
    SecretKey string `json:"secret_key"`
    Domain string `json:"domain"`
    PrivateKey string `json:"private_key"`
    PublicKey string `json:"public_key"`
    Region string `json:"region" default:"ap-guangzhou"`
}

func updateCredential(ctx context.Context, event DefineEvent) (string, error) {

    sync := Sync{
        SecretId:       issue.SecretId,
        SecretKey:      issue.SecretKey,
        Domain:         issue.CdnDomain,
        PrivateKeyData: privateKeyData,
        PublicKeyData:  publicKeyData,
        Region:         issue.Region,
    }

    return fmt.Sprintf("Hello %s!", event.Key1), nil
}

func main() {
    // Make the handler available for Remote Procedure Call by Cloud Function
    cloudfunction.Start(updateCredential)
}