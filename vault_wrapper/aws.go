package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    vault "github.com/hashicorp/vault/api"
    auth "github.com/hashicorp/vault/api/auth/aws"
)

func getSecretWithAWSAuthIAM() (string, error) {
    var (
        secret string
        value string
    )
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
    path, githubOutputVar, keyName := readSecretData()
    
    if keyName != "" {
        fmt.Println("printing from if statement")
        secret, err := client.KVv2(path).Get(context.Background(), keyName)
        fmt.Println("printing the secret - %v" , secret)
        value, ok := secret.Data["password"].(string)
        if !ok {
            return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["password"], secret.Data["password"])
        }
    } else {
        secret := client.KVv2(path).Get()
    } 
    
    
    // if err != nil {
    //     return "", fmt.Errorf("unable to read secret: %w", err)
    // }

    // data map can contain more than one key-value pair,
    // in this case we're just grabbing one of them
    

    fmt.Println("printing the secret - %v" , secret)

    return value, nil
}

func main() {
    if len(os.Args) < 4 {
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

func readSecretData() (string, string, string){
	secretData := "secrets/dev/kvv2/example foo | MY_PASSWORD"

    var (
        secretsPath string
        variable string
        key string
    )

	// Splitting the string by ' '
	parts := strings.Split(secretData, " ")
	if len(parts) >= 4 {
		secretsPath := parts[0]
		fmt.Println("Secrets Directory:", secretsPath)
	} else {
        secretsPath := parts[0]
        fmt.Println("Secrets Directory:", secretsPath)
    }

	// Splitting the string by ' | '
	variables := strings.Split(secretData, " | ")
	if len(variables) >= 2 {
		variable := variables[1]
		fmt.Println("Password:", variable)
	}

	// Splitting the second part by space to get 'foo'
	secondPart := strings.TrimSpace(variables[0])
	secondPartParts := strings.Split(secondPart, " ")
	if len(secondPartParts) >= 2 {
		key := secondPartParts[1]
		fmt.Println("Foo:", key)
	}

    return secretsPath, key, variable 
}
