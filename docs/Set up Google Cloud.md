# How to set up Direct Workload Identity Federation in Google Cloud

Source: https://github.com/google-github-actions/auth?tab=readme-ov-file#preferred-direct-workload-identity-federation


Note: The terminal commands below are applied to Command Prompt

1. Authorize to gcloud

```
gcloud auth login
```


1. Set variables for future use

```
set PROJECT_ID=<my-project-id>
```

```
set GITHUB_ORG=l122
```

check:

```
echo %PROJECT_ID%
```

```
echo %GITHUB_ORG%
```

1. (option) Create a service account

```
gcloud iam service-accounts create "ex-tracker-service-account" ^
    --project "%PROJECT_ID%"
```

save it

```
set GCP_SERVICE_ACCOUNT=<service-account@email>
```

check:

```
echo %GCP_SERVICE_ACCOUNT%
```

1. Create a Workload Identity Pool:

```
gcloud iam workload-identity-pools create "github" ^
  --project="%PROJECT_ID%" ^
  --location="global" ^
  --display-name="GitHub Actions Pool"
```

1. Get the full ID of the Workload Identity Pool and save it in a variable:

```
gcloud iam workload-identity-pools describe "github" ^
  --project="%PROJECT_ID%" ^
  --location="global" ^
  --format="value(name)"
```

save it:

```
set GCP_WORKLOAD_ID_PROVIDER=<output from above>
```

check:

```
echo %GCP_WORKLOAD_ID_PROVIDER%
```


1. Create a Workload Identity Provider in that pool:

```
gcloud iam workload-identity-pools providers create-oidc "my-repo" ^
  --project="%PROJECT_ID%" ^
  --location="global" ^
  --workload-identity-pool="github" ^
  --display-name="My GitHub repo Provider" ^
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner" ^
  --attribute-condition="assertion.repository_owner == '%GITHUB_ORG%" ^
  --issuer-uri="https://token.actions.githubusercontent.com"
```

1. Extract the Workload Identity Provider resource name:

```
gcloud iam workload-identity-pools providers describe "my-repo" ^
  --project="%PROJECT_ID%" ^
  --location="global" ^
  --workload-identity-pool="github" ^
  --format="value(name)"
```

save it

```
set GCP_WORKLOAD_ID_PROVIDER=<output from above>
```

check:

```
echo %GCP_WORKLOAD_ID_PROVIDER%
```

## Helpful gcloud commands

### Listing Service Accounts in a Project

```
gcloud iam service-accounts list --project=%PROJECT_ID%
```