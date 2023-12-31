package main

import (
    "context"
    "fmt"
    "os"
    "errors"
    "os/exec"
    "strings"

    vault "github.com/hashicorp/vault/api"
    auth "github.com/hashicorp/vault/api/auth/aws"
)

func getSecretWithAWSAuthIAM() (string, error) {
    role := os.Args[1]
    config := vault.DefaultConfig()

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
    parts := strings.Fields(secretData)
    // Grouping the elements into sets of four
	var groupedData [][]string
	for i := 0; i < len(parts); i += 4 {
		end := i + 4
		if end > len(parts) {
			end = len(parts)
		}
		groupedData = append(groupedData, parts[i:end])
	}
    for i, group := range groupedData {
		groupString := strings.Join(group, " ")
		fmt.Println("fetching input", i+1, groupString)
		mountPath, path, keyName, githubOutputVar := readSecretData(groupString)
        secret, err := client.KVv2(mountPath).Get(context.Background(), path)
        
        if secret == nil {
            fmt.Println("unable to fetch secrets from the path:", mountPath+path, "please check your path and the role")
            os.Exit(1)
        }
        
        value, ok := secret.Data[keyName].(string)
        if !ok {
            return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data[keyName], secret.Data[keyName])
        }

        var secretValue string
        if githubOutputVar != "" {
            secretValue = strings.TrimSpace(githubOutputVar) + "=" + value
        }else{
            secretValue = keyName + "=" + value
        }

        os.Setenv("secretValue", secretValue)
        commandToRun := fmt.Sprintf(`echo "$secretValue" >> "$GITHUB_OUTPUT"`)
        cmd := exec.Command("/bin/sh", "-c", commandToRun)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println("could not run command: ", err)
        }
        _ = out
	}
    return "Success", nil
}

func main() {
    if len(os.Args) < 4 {
        error := errors.New("Usage: ./aws.go $ROLE_NAME $VAULT_ADDR $VAULT_NAMESPACE")
        fmt.Println(error)
        os.Exit(1)
    }
    // setting env namespace, url
    os.Setenv("VAULT_ADDR", os.Args[3])
    os.Setenv("VAULT_NAMESPACE", os.Args[2])

    output, err := getSecretWithAWSAuthIAM()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    _ = output
}

func readSecretData(data string) (string, string, string, string){

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
	} else {
        secretsPath = parts[0]
        mountPathInput, secretsPath = formatSeretPath(secretsPath)
    }

	// Splitting the string by ' | '
	variables := strings.Split(data, " | ")
	if len(variables) >= 2 {
		variable = variables[1]
	}

	// Splitting the second part by space to get 'key'
	secondPart := strings.TrimSpace(variables[0])
	secondPartParts := strings.Split(secondPart, " ")
	if len(secondPartParts) >= 2 {
		key = secondPartParts[1]
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