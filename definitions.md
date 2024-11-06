Reduce your pods replicas to the amount defined on the following annotations:

```kube-pods-vacations/desired-replicas: number of desired pods```
```kube-pods-vacations/minimal-replicas: number of pods to be reduced```
```kube-pods-vacations/reduce-replicas-time-cron: cron format to reduce the pods```
```kube-pods-vacations/return-replicas-time-cron: cron format to return the desired amount of pods```
```kube-pods-vacations/last-execution-status: last time the kube-pods-vacations scheduler was executed```

