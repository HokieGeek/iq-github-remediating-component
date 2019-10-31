# iq-merge-review-remediations [![DepShield Badge](https://depshield.sonatype.org/badges/sonatype-nexus-community/iq-merge-review-remediations/depshield.svg)](https://depshield.github.io)

AWS Lambda which uses your Sonatype Nexus IQ instance to capture GitHub Pull Requests and/or GitLab Merge Requests from your repos and adds inline comments with suggestions on versions to upgrade your vulnerable open source components.

## How to use

1. Build and upload as AWS Lambda
2. Add your webhook to your repo's config with the following payload URL:
`<LAMBDA_API_GATEWAY_ENDPOINT>?iq_url=<IQ_SERVER_PORT>&iq_auth=<IQ_USER>:<IQ_PASS>&iq_app=<IQ_APP>&token=<ACCESS_TOKEN>`

## Supported languages
* Java (maven, gradle)
* go (go modules)
* C# / .net (nuget)
* Javascript / Typescript (npm)
* Ruby (rubygems)

## Examples

### GitHub Pull Request
https://github.com/HokieGeek/various-manifests/pull/49/files

### GitLab Merge Request
https://gitlab.com/HokieGeek/various-manifests/merge_requests/5/diffs