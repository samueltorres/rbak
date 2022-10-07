#!/bin/bash

cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "8760h"
      }
    }
  }
}
EOF

cat > ca-csr.json <<EOF
{
  "CN": "Kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Portland",
      "O": "Kubernetes",
      "OU": "CA",
      "ST": "Oregon"
    }
  ]
}
EOF

cfssl gencert -initca ca-csr.json | cfssljson -bare ca

cat > server-csr.json <<EOF
{
  "CN": "admission",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Portland",
      "O": "Kubernetes",
      "OU": "Kubernetes",
      "ST": "Oregon"
    }
  ]
}
EOF

cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -hostname=rbak.default.svc \
  -profile=server \
  server-csr.json | cfssljson -bare server


echo "Generating secret"
kubectl create secret tls rbak-webhook-tls \
  --cert=server.pem \
  --key=server-key.pem \
  --dry-run=client -o yaml \
  > ./hack/webhook-secret.yaml

echo "Generating helm values with ca bundle"
CA_BUNDLE=$(cat ca.pem | base64 - | tr -d '\n' )
cat <<EOF >./hack/webhook-ca-values.yaml
webhook:
  cert:
    enabled: true
    caBundle: $CA_BUNDLE
EOF


echo "Removing non-needed files"
rm ca.* ca-* server.* server-*