We need to fill this secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ssh-key
  annotations:
    tekton.dev/git-0: github.com # Described below
type: kubernetes.io/ssh-auth
stringData:
  ssh-privatekey: <private-key>
  # This is non-standard, but its use is encouraged to make this more secure.
  # If it is not provided then the git server's public key will be requested
  # when the repo is first fetched.
  known_hosts: <known-hosts>
```

Then we need to add the secret to the service account.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: build-bot
secrets:
  - name: ssh-key
```