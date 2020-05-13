# Apache Benchmark Operator

An Kubernetes operator for running Jobs using the Apache HTTP server benchmarking tool [ab][ab_docs].

## Installation

The following steps will install the operator into a Kubernetes cluster. Clone this repository somewhere locally to get
started.

NOTE: You will need the `cluster-admin` or equivalent role on the cluster to perform several of these steps.

Create a namespace as needed for your environment.

``` bash
kubectl create ns benchmark
```

Create the RBAC resources for the operator.

``` bash
kubectl apply -n benchmark -f deploy/service_account.yaml
kubectl apply -n benchmark -f deploy/role.yaml
kubectl apply -n benchmark -f deploy/role_binding.yaml
```

Create the CRD that is managed by the operator.

``` bash
kubectl apply -f deploy/crds/httpd.apache.org_apachebenches_crd.yaml
```

Create the Deployment to run the operator.

``` bash
kubectl apply -n benchmark -f deploy/operator.yaml
```

If these steps all complete successfully, the operator Pod should be running in the desired namespace.

``` bash
kubectl get pods -n benchmark
```

``` bash
NAME                                     READY   STATUS    RESTARTS   AGE
apache-bench-operator-58dfcbd5fd-ln6v9   1/1     Running   0          13s
```

## Usage

With the operator running, create a new `ApacheBench` resource that defines the desired state for the benchmark
process.

``` yaml
apiVersion: httpd.apache.org/v1alpha1
kind: ApacheBench
metadata:
  name: example-apache-bench
  labels:
    example: basic
spec:
  url: http://httpd.apache.org/
```

This example is available in the `docs/examples` directory.

``` bash
kubectl apply -n benchmark -f docs/examples/apachebench-basic.yaml
```

``` bash
apachebench.httpd.apache.org/example-apache-bench created
```

View the `ApacheBench` resources.

``` bash
kubectl get ab -n benchmark
```

``` bash
NAME                   AGE
example-apache-bench   106s
```

A Job will be created to run the benchmark.

``` bash
kubectl get jobs -n benchmark
```

``` bash
NAME                   COMPLETIONS   DURATION   AGE
example-apache-bench   1/1           2s         2m19s
```

The Job will create a single Pod and should not take very long to run.

``` bash
kubectl get pods -n benchmark
```

``` bash
NAME                                     READY   STATUS      RESTARTS   AGE
apache-bench-operator-58dfcbd5fd-ln6v9   1/1     Running     0          22m
example-apache-bench-wwx7h               0/1     Completed   0          4m8s
```

View the logs for the Pod to get the result of the benchmark process.

``` bash
kubectl logs -n benchmark example-apache-bench-wwx7h
```

``` bash
This is ApacheBench, Version 2.3 <$Revision: 1874286 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking httpd.apache.org (be patient).....done


Server Software:        Apache/2.4.18
Server Hostname:        httpd.apache.org
Server Port:            80

Document Path:          /
Document Length:        9696 bytes

Concurrency Level:      1
Time taken for tests:   0.013 seconds
Complete requests:      1
Failed requests:        0
Total transferred:      9969 bytes
HTML transferred:       9696 bytes
Requests per second:    79.33 [#/sec] (mean)
Time per request:       12.605 [ms] (mean)
Time per request:       12.605 [ms] (mean, across all concurrent requests)
Transfer rate:          772.34 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        6    6   0.0      6       6
Processing:     6    6   0.0      6       6
Waiting:        6    6   0.0      6       6
Total:         13   13   0.0     13      13
```

## License

The Apache Benchmark Operator is released under the Apache 2.0 license. See the [LICENSE][license_file] file for details.

[ab_docs]:https://httpd.apache.org/docs/2.4/programs/ab.html
[license_file]:./LICENSE
