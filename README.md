# operator-nodejs

POC of "day-2" operations for Node.js, equivalent to the
[OpenLiberty Day-2 Operations][].

[OpenLiberty Day-2 Operations]: https://github.com/OpenLiberty/open-liberty-operator/blob/master/doc/user-guide.md#request-server-dump

## Try it out

1. Deploy a report-capable application with appsody.
2. Use operator-nodejs to trigger a report.

### Deploy an appsody app built with a report-capable stack:

1. Package the stack that supports diagnostic report
   1. Checkout https://github.com/appsody/stacks/pull/722
   2. cd stacks/incubator/nodejs-expres
   3. appsody stack package
2. Build an app based on the stack
   1. mkdir ex-nodejs-report
   2. cd ex-nodejs-report
   3. appsody init dev.local/nodejs-express
   4. appsody build
3. Confirm it runs locally and generates reports
   1. docker run -it dev.local/ex-nodejs-report
   2. (in another terminal) docker exec CONTAINER /bin/sh -c 'kill -USR2 $(pidof node)'
   3. ... see json report logged to stdout... note it is not (yet) consumable by
      a JSON log system because it is multi-line:
      - https://github.com/nodejs/node/pull/32254 (landed)
      - https://github.com/nodejs/node/pull/32497 (WIP)
4. Deploy app with appsody
   1. appsody deploy -t octet/ex-nodejs-report --push
   2. kubectl get pod -L nodejs.appsody.dev/report
   3. Remember the name of the reportable nodejs pod for later


### Use operator-nodejs to trigger a report.

In separate terminals:
- run operator: make kube-init local-run
- show logs: kubectl logs -f pod/PODNAME-FROM-ABOVE
- trigger report: make example-report

Reporting is done by creating a NodejsReport custom resource:
```
apiVersion: opnodejs.appsody.dev/v1beta1
kind: NodejsReport
metadata:
  name: example-nodejsreport
spec:
  podName: REPLACEME_PODNAME
```

## TODO

- [x] rename "diagnosticreport" to "report", octetcloud to appsody.dev
- feature: report
  - [x] send signal
  - [x] check label for support (label searchable)
  - [x] appsody stack with report integration
  - [ ] check annotation for signal name?
  - [ ] when report is deleted, there is a log msg about being reconciled, but
        nothing saying why nothing is happening
  - [ ] conditions history?
    - try? https://github.com/redhat-cop/operator-utils
  - [ ] kAppNav integration
  - [ ] appsody stack: PR?
  XXX does it report on handled promises? when exit code is not 0?
  - [ ] tests?
    - https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md

- [ ] https://github.ibm.com/runtimes/squad-node/wiki/Playbacks

- [ ] add lint, etc, to operator from shorty

- recreate with kubebuilder?
  - https://github.com/kubernetes-sigs/kubebuilder
    - better than op-sdk? should try, looks promising
    - http://banjiewen.net/operator-sdk.html: use kubebuilder...
> etcd-operator and prometheus-operator are implemented with basic client-go
> based controllers. No Operator SDK, no Kubebuilder, etc.

- feature: suport label of `*`
  - [ ] ...
  - build a queue of work, and a goroutine to do the work, fetch crd at
    beginning of each work to make sure its still current, somehow mark status
    per target. deside how parallel to make the signalling... or don't care?
  - https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#-strong-read-operations-pod-v1-core-strong-
  - Query:
	// Find all reportable pods
	list := &corev1.PodList{}
	// opts := client.MatchingLabels{"nodejs.appsody.dev/reportable": "true"}
	opts := client.MatchingLabels{"run": "ex"}
	err = mgr.GetAPIReader().List(context.TODO(), list, opts)
	if err != nil {
		log.Error(err, "get reportable")
		return nil
	}
	for _, pod := range list.Items {
		log.Info("reportable", "pod", pod)
	}
	panic("x")

- feature: debug port (has a state, so fits into kube better)

- feature: dump heap profile/cpu profile?

- feature: `kubectl nodejs-report` should be possible with an AA
  - Aggregated API Servers (AA)
  - https://medium.com/@cloudark/kubernetes-operators-when-how-and-the-gotchas-to-keep-in-mind-b13be9830346
