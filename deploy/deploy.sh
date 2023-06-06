#!/usr/bin/env sh

set -eux

# Create the namespace
kubectl create namespace $NAMESPACE || true

# Delete Kubernetes secret if exists
kubectl delete secret clinical-service-account --namespace $NAMESPACE || true

# Create GCP service account file
cat $GOOGLE_APPLICATION_CREDENTIALS >> ./service-account.json

# Recreate service account file as Kubernetes secret
kubectl create secret generic clinical-service-account \
    --namespace $NAMESPACE \
    --from-file=key.json=./service-account.json

helm upgrade \
    --install \
    --debug \
    --create-namespace \
    --namespace "${NAMESPACE}" \
    --set app.replicaCount="${APP_REPLICA_COUNT}" \
    --set service.port="${PORT}"\
    --set app.container.image="${DOCKER_IMAGE_TAG}"\
    --set app.container.env.googleCloudProject="${GOOGLE_CLOUD_PROJECT}"\
    --set app.container.env.environment="${ENVIRONMENT}"\
    --set app.container.env.googleProjectNumber="${GOOGLE_PROJECT_NUMBER}"\
    --set app.container.env.sentryDSN="${SENTRY_DSN}"\
    --set app.container.env.cloudHealthPubsubTopic="${CLOUD_HEALTH_PUBSUB_TOPIC}"\
    --set app.container.env.cloudHealthDatasetID="${CLOUD_HEALTH_DATASET_ID}"\
    --set app.container.env.cloudHealthDatasetLocation="${CLOUD_HEALTH_DATASET_LOCATION}"\
    --set app.container.env.cloudHealthFHIRStoreID="${CLOUD_HEALTH_FHIRSTORE_ID}"\
    --set app.container.env.openConceptLabToken="${OPENCONCEPTLAB_TOKEN}"\
    --set app.container.env.serviceHost="${SERVICE_HOST}"\
    --set app.container.env.openConceptAPIUrl="${OPENCONCEPTLAB_API_URL}"\
    --set app.container.env.savannahAdminEmail="${SAVANNAH_ADMIN_EMAIL}"\
    --set app.container.env.authserverEndpoint="${AUTHSERVER_ENDPOINT}"\
    --set app.container.env.clientID="${CLIENT_ID}"\
    --set app.container.env.clientSecret="${CLIENT_SECRET}"\
    --set app.container.env.authUsername="${AUTH_USERNAME}"\
    --set app.container.env.authPassword="${AUTH_PASSWORD}"\
    --set app.container.env.grantType="${GRANT_TYPE}"\
    --set app.container.env.mycarehubClientID="${MYCAREHUB_CLIENT_ID}"\
    --set app.container.env.mycarehubClientSecret="${MYCAREHUB_CLIENT_SECRET}"\
    --set app.container.env.mycarehubIntrospectURL="${MYCAREHUB_INTROSPECT_URL}"\
    --set app.container.env.clinicalBucketName="${CLINICAL_BUCKET_NAME}"\
    --set networking.issuer.name="letsencrypt-prod"\
    --set app.container.env.jwtKey="${JWT_KEY}"\
    --set networking.issuer.privateKeySecretRef="letsencrypt-prod"\
    --set networking.ingress.host="${APP_DOMAIN}"\
    --wait \
    --timeout 300s \
    -f ./charts/clinical/values.yaml \
    $APP_NAME \
    ./charts/clinical