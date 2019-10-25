# iq-github-remediating-component

AWS Lambda which uses your IQ instance to capture GitHub pull requests from your repos and adds comments to the PR with suggestions on versions to upgrade your vulnerable open source components.

Add your webhook with the following Payload URL:
`<LAMBDA_API_GATEWAY_ENDPOINT>?iq_url=<IQ_SERVER_PORT>&iq_auth=<IQ_USER>:<IQ_PASS>&iq_app=<IQ_APP>&token=<GITHUB_TOKEN>`
