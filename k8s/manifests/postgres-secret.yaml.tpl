---
# Source: postgresql/templates/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: esl-games-postgresql
  labels:
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-10.2.0
    app.kubernetes.io/instance: esl-games
    app.kubernetes.io/managed-by: Helm
type: Opaque
data:
  postgresql-postgres-password: {{.Password1}}
  postgresql-password: {{.Password2}}
