# TODO

This file contains task that are being work on or will be work on
shortly.

If you would like to help with the development of AEGo, please follow
these steps:

1. Select an item from the list that you would like to work on.
2. Search the the [issues track](https://github.com/scotch/aego/issues/) to see if someone has begun work on the issue.
3. If the issue has not been started, create a [new issue](https://github.com/scotch/aego/issues/new) outlining the problem as you understand it. 
4. Create a fork of AEGo the project.
5. Solve the problem.
6. Use [gofmt](http://golang.org/cmd/gofmt/) to format your code prior to commit.
7. If this is your first time contributing add your name to CONTRIBUTERS.md and AUTHORS.md.
8. Create a pull request.

### v1/auth

### v1/auth/password
- Make use of [subtle.ConstantTimeCompare](http://golang.org/pkg/crypto/subtle/#ConstantTimeCompare) to prevent timing attacks
- Add salt to bcrypt
  - salt should come from config setting. A proper solution for loading
    configs prior to app initialization is still still in the works, therefore this will have to wait.
 
### v1/auth/google
### v1/auth/facebook
### v1/auth/github
