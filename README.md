# Argo Workflows Gene
Argo Workflows Gene enables scientific workflows written in WDL and CWL languages to run on Kubernetes.

# Overview

Argo Workflows Gene allows scientific workflows to run seamlessly on Kubernetes. 
It relies on the CNCF graduation project Argo Workflows and supports WDL and CWL languages. 
It allows scientific computing to enjoy many advantages of cloud native, such as scalability and observability.

# What is Argo Workflows Gene


# Why Argo Workflows Gene Needed

Kubernetes has become a cloud-native operating system and the standard for managing applications.
Scientific computing workflows also seek to be managed on Kubernetes so that they can enjoy the advantages brought by the ecosystem.
However, they are often written using WDL and CWL languages, making them difficult to run on Kubernetes.

Now running these workflows on Kubernetes usually uses the Task Execution Service (TES) specification proposed by GA4GH.
The disadvantage is that it is not convenient to maintain and the cost is high (one task corresponds to 4 pods).

So we need a project to make scientific computing workflows seamless on Kubernetes. 

# How to use

client mode

```
argogene submit wdl --file helloworld.wdl --parameter input.json
```


