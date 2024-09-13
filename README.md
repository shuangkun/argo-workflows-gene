# Argo Workflows Gene

Argo Workflows Gene enables scientific workflows to run on Kubernets, and it supports WDL and CWL languages.

# Overview

Argo Workflows Gene allows scientific workflows to run seamlessly on Kubernetes. 
It relies on the CNCF graduation project Argo Workflows and supports WDL and CWL languages. 
It allows scientific computing to enjoy many advantages of cloud native, such as scalability and observability.

# Why Argo Workflows Gene Needed

Kubernetes has become a cloud-native operating system and the standard for managing applications.
Scientific computing workflows also seek to be managed on Kubernetes so that they can enjoy the advantages brought by the ecosystem.
But now there are many difficulties in running these workflows written in WDL and CWL languages on kubernetes.

* Difficult to use and maintain

 Now running these workflows on Kubernetes usually uses the Task Execution Service (TES) specification proposed by GA4GH.
The disadvantage is that it is not convenient and difficult to maintain.

* Cost is high.

A task is split into four. More tasks are started and more resources are consumed.

* Low performance

When some large-scale data exchange is required, performance cannot be guaranteed, and only local and shared storage is supported.

* Lack of advanced capabilities

Some batch computing scheduling capabilities in Kubernetes, such as Gang and topology awareness, are difficult to use.

These challenges make it very difficult for the scientific computing industry to migrate to the Kubernetes ecosystem and use Kubernetes. 
Argo Workflows Gene is designed to help scientific computing workflows run seamlessly on Kubernetes based on Argo Workflow, the powerful engine of Kubernetes.

# Start

```
git clone git@github.com:shuangkun/argo-workflows-gene.git
sh build.sh
```

# Mode

client mode

```
argogene submit wdl --file helloworld.wdl --parameter input.json
```




