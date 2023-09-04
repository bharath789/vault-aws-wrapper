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

    
    
    secretData := os.Args[4]
    fmt.Println("print secret data - %v" ,secretData)
    // add the secret logic fetch multiple secrets
    mountPath, path, keyName, githubOutputVar := readSecretData(secretData)
    
    secret, err := client.KVv2(mountPath).Get(context.Background(), path)
    fmt.Printf("printing the secret: %v\n" , secret)
    fmt.Printf("printng the mountPath: %v\n", mountPath)
    fmt.Printf("printng the path: %v\n", path)
    fmt.Printf("printing githubOutputVar: %v\n", githubOutputVar)
    fmt.Printf("printing keyName: %v\n", keyName)
    fmt.Printf("printing secret Data - %#v\n", secret.Data[keyName])
    value, ok := secret.Data[keyName].(string)
    if !ok {
        return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data[keyName], secret.Data[keyName])
    }
    fmt.Println("Secret Value:", value)
    secretValue := keyName + "=" + value
    fmt.Println("printing secretVaule : ", secretValue)
    // fmt.Println(fmt.Sprintf(`::set-output name=%s::%s`, keyName, value))
    os.Setenv("GITHUB_OUTPUT", secretValue)

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

    output, err := getSecretWithAWSAuthIAM()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf(output) 
}

func readSecretData(data string) (string, string, string, string){
	// data := "dev/kvv2/example foo | MY_PASSWORD"

    var (
        secretsPath string
        mountPathInput string
        variable string
        key string
    )
	// Splitting the string by ' '
	parts := strings.Split(data, " ")
	if len(parts) >= 4 {
        secretsPath = parts[0]
        mountPathInput, secretsPath = formatSeretPath(secretsPath)
        fmt.Println("Secrets Directory:", mountPathInput)	
	} else {
        secretsPath = parts[0]
        mountPathInput, secretsPath = formatSeretPath(secretsPath)
        fmt.Println("Secrets Directory:", mountPathInput)
    }

	// Splitting the string by ' | '
	variables := strings.Split(data, " | ")
	if len(variables) >= 2 {
		variable = variables[1]
		fmt.Println("github variable:", variable)
	}

	// Splitting the second part by space to get 'key'
	secondPart := strings.TrimSpace(variables[0])
	secondPartParts := strings.Split(secondPart, " ")
	if len(secondPartParts) >= 2 {
		key = secondPartParts[1]
		fmt.Println("key:", key)
	}

    return mountPathInput, secretsPath, key, variable 
}

func formatSeretPath(fmtPath string) (string, string){
    var (
        fmtPathData string
        lastPosition int
    )
    secretsPathNew := strings.Split(fmtPath, "/")
    
    for i :=0; i < (len(secretsPathNew) - 1); i++ {
        fmtPathData += secretsPathNew[i] + "/"
        lastPosition = i + 1
    }
    return fmtPathData, secretsPathNew[lastPosition]
}