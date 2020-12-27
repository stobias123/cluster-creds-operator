# Github Credential Operator.

I hate looking up how to create creds for `docker pull` and  `github` in kubernetes every time I set up a new cluster.

It also be better practice to use new secrets per namespace or account.. that doesn't always happen.

So overall, 
1. Create ephemeral credentials for various services (gcr.io,docker,github.com,gitlab.com)
2. Attach properly (and automatically) to proper service accounts across targeted namespaces across the cluster.