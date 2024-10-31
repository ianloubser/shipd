## go deploy

This is a nifty lightweight server I use to deploy updates and rollout changes to my side hustles. Docker is overkill, and likewise any AWS Fargate or other offerings. I run this API as a service and it takes care of pulling a git repo, doing setup & teardown and restart of service to update. 
This means my repo workflow/pipeline can just do a POST request to this API with commit SHA and that version will get deployed. 

### VPS Setup
I setup base bitnami Ubuntu image wherever I rent my VPS from. Then next step is to install `pm2` globally to manage all service lifecycles to be run (even this deploy server is run with pm2).

And then I use global install of `nvm` so different projects can use different node versions, and lastly a global install of `Caddyserver`. This brings my global installs to 3 command line arguments to get set up.

The only next step to get going is then to create the `.toml` file to configure each project, and thats it!.

All of the above is also automated, and you could just run the `vm-install.sh` script to get it all set up, example:
```bash
curl -o- https://raw.githubusercontent.com/ianloubser/smol-deploy/vm-install.sh {wildcard-domain-here} | bash
```

You'd only need to SSH into VPS now for any normal maintenance issues but your general deploy lifecycle should be sorted :) 

For a breakdown of how the service works, have a read here [https://ianloubser.github.io/posts/deploy-your-sideprojects/](https://ianloubser.github.io/posts/deploy-your-sideprojects/)