apiVersion: apps/v1
kind: Deployment
metadata:
  name: lesson-controller-manager
  namespace: system
spec:
  template:
    spec:
      containers:
      - name: manager
        args:
        - --gameserver-image LATEST_IMAGE