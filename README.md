# todo

## Usage

```
todo -vvvv
todo set -l cluster=ar1 -l role=worker ar1w0101.ncc ar1w0102.ncc
todo get -e "cluster==ar1"
todo get -e "cluster==ar1&&role=worker"
todo get -e "cluster==ar1" -o wide
todo get -e "cluster==ar1" -o yaml
todo run -e "cluster==ar1&&role=worker" -t ssh -- sudo journalctl
todo run -e "cluster==ar1&&role=worker" -t ssh -w 10 -- sudo journalctl
todo del ar1w0101.ncc ar1w0102.ncc

todo template ssh 'ssh {{.user}}@{{.name}} -C {{args}}'
```


