---
title: "LIDC Experiment on AWS"
weight: -60
---

# Running the LIDC experiment on the AWS cluster

Start up the cluster

```
# Startup a cluster called 'lidc' using the 'lidccluster' template
. venv/bin/activate
starcluster start -c lidccluster lidc
```


Copy software the slow way...

```
starcluster put lidc --user sgeadmin ClusterSoftware/ /software/
starcluster put lidc --user sgeadmin  devops/lidc.sh /home/sgeadmin/lidc.sh
```

or use rsync

```
devops/cluster-install.sh
```

Use the [elastic load balancer](http://star.mit.edu/cluster/docs/latest/manual/load_balancer.html) to ramp up nodes.

```
# -m 10     -- maximum number of nodes
# -a 2      -- add 2 nodes at a time
# -p -P dir -- save a series of plots to dir
# -d -D dir -- save a CSV of stats to dir
mkdir cluster_stats
starcluster loadbalance -p -P cluster_stats -d -D cluster_stats lidc
```


### Fixing problems

