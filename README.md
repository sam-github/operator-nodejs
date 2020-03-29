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
