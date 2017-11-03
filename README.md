# S3 Artifactory

The idea of this tool is to allow querying and fetching artifacts between different
CD (Continous Delivery) systems that share nothing but the deployment artifacts.

- The CI job will build and publish the artifact
- The Staging deployment will fetch the latest build artifact and deploy it to staging
- The QA deployment job will fetch the latest staging artifact
- The Production deployment job will fetch the latest promoted QA build

All CI/Staging/QA/Production can be totally decoupled on running on different data-centers.

This can run as a CLI tool with the desired argumentor can be developed as a HTTP service
that exposes the artifacts information on REST interface.
