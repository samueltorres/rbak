# rbak

rbak is a kubernetes rbac analyser that allows you to understand what roles service accounts, users and groups need based on their interactions with the kube-apiserver.

## Getting Started

`rbak` works as a `ValidatingWebhook` that records every interaction with the kube-apiserver and then builds a report of the resources and verbs a given service account, user or group used so far. This is really useful to ensure you're following the least-privilege principle, because you're able to compare the access given vs the access needed.

### Running on the cluster

1. Generate Webhook Certificates

```sh
make gen-certs
```

2. Install it using the Helm Chart:

```sh
make install
```


## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

