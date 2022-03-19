package main

import (
    "context"
    "fmt"
    "github.com/actors315/incubator-certificate-ticket/qcloud-tools/certificate"
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

    client := certificate.LBSync{
        Sync: certificate.Sync{
            SecretId:       event.SecretId,
            SecretKey:      event.SecretKey,
            Domain:         event.Domain,
            PrivateKeyData: event.PrivateKey,
            PublicKeyData:  event.PublicKey,
            Region:         event.Region,
        },
    }

    result := client.UpdateCredential()
    if !result {
        errMsg := fmt.Sprintf("update certificate failed, %s \n" , event.Domain)
        return errMsg, fmt.Errorf(errMsg)
    }

    return fmt.Sprintf("update certificate success, %s \n", event.Domain), nil
}

func main() {
    // Make the handler available for Remote Procedure Call by Cloud Function
    cloudfunction.Start(updateCredential)
}