# Knative Gitlab Source-to-app with push

This PoC was created during: https://www.hack8s.de
-> Create Github Repo
-> New push should trigger a new serving revision

## Goal

Trigger a build with Knative Serving [source-to-url](https://github.com/knative/docs/tree/master/serving/samples/source-to-url-go) when a new commit is pushed to the remote repository (normally the image would only be build once during the creationg of the serving).

## Notes

Fun with knative:

- https://github.com/knative/docs/issues/817
- https://github.com/knative/build/issues/566

Fix for 566:

```bash
# use 0.3.0
kubectl apply -f https://github.com/knative/build/releases/download/v0.3.0/release.yaml
```

## Prequisite

- K8s cluster with knative: (https://github.com/knative/docs/tree/master/install)

## Install kaniko for building

```bash
# I used the commit 58f31c46fc2009a766fc2c09d2f340b790816a64 maybe a future commit may break everything
kubectl apply -f https://raw.githubusercontent.com/knative/build-templates/master/kaniko/kaniko.yaml
```

## Set up DNS

First we fetch the ingress ip:

```bash
export $INGRESSGATEWAY=istio-ingressgateway
ING_IP=$(kubectl get svc $INGRESSGATEWAY --namespace istio-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')
echo "$ING_IP.nip.io"
```

Adjust the `config-domain` default domain with the domain above (just write the domain on a top-level: https://github.com/knative/docs/blob/master/serving/using-a-custom-domain.md):

```bash
kubectl -n knative-serving edit cm config-domain
```

## Prepare the build

In order to make it possible that kaniko can push images to a Docker repository we need to create a Kubernetes secret containing the secrets that kaniko uses to push the docker image to the docker hub (or another repoisotry).

echo -ne "$user" | base64

```bash
echo -ne "$user" | base64
echo -ne "$password" | base64
cp build-sa.yaml.tpml build-sa.yaml
```

Now we can apply it:

```bash
kubectl apply -f build-sa.yaml
```

Ensure that all parts are working:

```bash
kubectl apply -f test-build.yaml
```

Wait and check if the build succeeded:

```bash
kubectl get build
```

Now we clean up:

```bash
kubectl delete -f test-build.yaml
```

Now you need to adjust in the `demo-app.yaml` file the Docker registry (now you can apply it):

```bash
kubectl apply -f demo-app.yaml
```

Ensure that this triggers a build:

```bash
kubectl get build
```

wait until the build is done and the ksvc get's a domain:

```bash
kubectl get ksvc
```

And now we can test the app:

```bash
curl app-from-source.default.35.204.237.123.nip.io
...
```

## Build new revision on push

Create new ServiceAccount with requiered roles and deploy the event trigger (note this endpoint is public available! You probably don't want this):

```bash
kubectl apply -f event-trigger.yaml
```

Change the personal access token in `github-source.yaml` and run:

```bash
kubectl apply -f github-source.yaml
```

## Have fun

Create a new push and wait for the new build:

```bash
kubectl get build
```

Warning this will trigger a new build for any push on this repo (even if it's not master and the Docker image will be build from the master branch)
