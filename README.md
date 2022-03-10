# pw

```
pw template gcloud-ssh -- gcloud compute --project "{{.project}}" ssh --zone "{{.zone}}" "{{.name}}" -- "{{args}}"
pw template gcloud-scp -- gcloud compute --project '{{.project}}' scp --zone "{{.zone}}" "{{arg 1}}" {{.name}}:"{{arg 2}}"
seq -f "k8s-node-%03g" 1 32 | xargs pw set -l zone=us-west1-a -l project=pw1
seq -f "k8s-node-%03g" 1 32 | xargs pw set -l zone=us-west1-b -l project=pw2
seq -f "k8s-node-%03g" 1 32 | xargs pw set -l zone=us-west1-a -l project=pw1
seq -f "k8s-node-%03g" 1 32 | xargs pw set -l zone=us-west1-a -l project=pw2

pw run -t gcloud-ssh -- echo hello
pw run -t gcloud-scp -- test.sh t.sh
pw run -e 'project =="pw1" && name endsWith "1"' -t gcloud-ssh -- echo hello
```

## Usage

### Set

### Get

### Del

### Run

### Template

### Shortcut Get

### Shortcut Set

```
pw -vvvv
pw set -l cluster=ar1 -l role=worker ar1w0101.ncc ar1w0102.ncc
pw get -e "cluster==ar1"
pw get -e "cluster==ar1&&role=worker"
pw get -e "cluster==ar1" -o wide
pw get -e "cluster==ar1" -o yaml
pw run -e "cluster==ar1&&role=worker" -t ssh -- sudo journalctl
pw run -e "cluster==ar1&&role=worker" -t ssh -w 10 -- sudo journalctl
pw del ar1w0101.ncc ar1w0102.ncc

pw template ssh 'ssh {{.user}}@{{.name}} -C {{args}}'
```
