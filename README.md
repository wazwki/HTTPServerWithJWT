# How to create golang HTTP-server with JWT-authentication

## Step 1: 
Install library for use jwt in project:
```bash 
go get "github.com/golang-jwt/jwt/v5"
```
```go
import (
	"github.com/golang-jwt/jwt/v5"
)
```


## Step 2: 

Create made-token and check-token functions:

```go
func MadeToken(username, email string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
    jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil { 
        return "", err 
    }

	return tokenString, nil
}
```

```go
func CheckToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, 
        func(token *jwt.Token) (interface{}, error) {
		    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			    return nil, fmt.Errorf("unexpected: %v", token.Header["alg"])
		    }
		    return secretKey, nil
	    }
    )
	if err != nil { 
        return nil, err 
    }

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
```


## Step 3: 

Create token in login handler like this:

```go
secretKey := []byte("your-secret-key")
token, _ := jwt.MadeToken(User.Username, User.Email, secretKey)
```

Check token in auth-only handler like this:

```go
secretKey := []byte("your-secret-key")
tokenString := r.Header.Get("Authorization")
claims, _ := jwt.CheckToken(tokenString, secretKey)

username, ok := claims["username"].(string)

email, ok := claims["email"].(string)
```
