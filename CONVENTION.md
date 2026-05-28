### Code Convention

- filename and package name should use snake_case convention
- variabels name and function name should use camelCase convention
- contracts defined in [```domain```](./domain/) should use PascalCase convention
- don't use 1 letter variable name
- all structs that are an implementation of contracts in [```internal```](./internal/) directory should be private, thus making the struct name follows camelCase convention
- code's identation style is space with size 4
- when importing packages, prioritze standard library first, thirdy party library and finally local package.
example:
```go
import (
    "fmt"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/BangNopall/hology8-be/pkg/log"
)
```
- if a function accepts more than 3 parameters, the parameters should be expanding towards bottom, not right. 
example: 
```go
func Send(
    from string, 
    to string,
    title string,
    description string, 
    file File
) error {

}
```
- when handling an error that was occured by the application, always logs it with [```pkg/log```](./pkg/log/log.go)
- logging an info or error must use the following convention
```go
// logging with package defined in pkg/log
log.Info(log.LogInfo{
    "data": data
}, "[File Name in All Caps without .go Extension Separated By Space][Method Name] message")

// example 
log.Error(log.LogInfo{
    "error": err.Error(),
}, "[USER REPOSITORY][FetchByEmail] failed to fetch by email")
```
- handle responses in controller with [```pkg/helpers/http/response```](./pkg/helpers/http/response/response.go)
- unit testing must test all the functions that was defined by contracts defined in [```domain```](./domain/)

### Commit Convention

All commits message must follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)

### API Naming Convention

For API Naming convention, you can read the following guide : [API Naming Convention](https://restfulapi.net/resource-naming/) and prefix all endpoints with ```/api/v1```

### Pushing Changes

After you commit all your changes, you should push your changes to your personal branch first. Your personal branch name should be your nickname. After that, you can send pull request to dev branch. Owner will merge dev branch to master branch when dev branch passes all unit testing and ready to be deployed on VPS