package main

import (
    "context"
    "fmt"
    "os"

    vault "github.com/hashicorp/vault/api"
    auth "github.com/hashicorp/vault/api/auth/aws"
)

func getSecretWithAWSAuthIAM() (string, error) {
    role := os.Args[1]
    config := vault.DefaultConfig() // modify for more granular configuration

    client, err := vault.NewClient(config)
    if err != nil {
        return "", fmt.Errorf("unable to initialize Vault client: %w", err)
    }

    awsAuth, err := auth.NewAWSAuth(
        auth.WithRole(role), 
    )
    if err != nil {
        return "", fmt.Errorf("unable to initialize AWS auth method: %w", err)
    }

    authInfo, err := client.Auth().Login(context.Background(), awsAuth)
    if err != nil {
        return "", fmt.Errorf("unable to login to AWS auth method: %w", err)
    }
    if authInfo == nil {
        return "", fmt.Errorf("no auth info was returned after login")
    }

    
    
    secretData := os.Args[5]
    fmt.Println("print secret data - %v" ,secretData)
    // add the secret logic fetch multiple secrets 
    secret, err := client.KVv2("secret").Get(context.Background(), "creds")
    if err != nil {
        return "", fmt.Errorf("unable to read secret: %w", err)
    }

    // data map can contain more than one key-value pair,
    // in this case we're just grabbing one of them
    value, ok := secret.Data["password"].(string)
    if !ok {
        return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["password"], secret.Data["password"])
    }

    return value, nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: ./aws.go $ROLE_NAME $VAULT_ADDR $VAULT_NAMESPACE")
        return
    }
    // setting env namespace, url
    os.Setenv("VAULT_ADDR", os.Args[3])
    os.Setenv("VAULT_NAMESPACE", os.Args[2])

    secretValue, err := getSecretWithAWSAuthIAM()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // fmt.Println("Secret Value:", secretValue)
    os.Setenv("GITHUB_OUTPUT", secretValue)
}