# Running the LIDC experiment on the AWS cluster

Start up the cluster

```bash
# Startup a cluster called 'lidc' using the 'lidccluster' template
. venv/bin/activate
starcluster start -c lidccluster lidc
```


Copy software the slow way...

```bash
starcluster put lidc --user sgeadmin ClusterSoftware/ /software/
starcluster put lidc --user sgeadmin  devops/lidc.sh /home/sgeadmin/lidc.sh
```

or use rsync

```bash
devops/cluster-install.sh
```

Log in to the cluster as `sgeadmin`:

```
starcluster sshmaster --user sgeadmin lidc
```

Launch the processing job:

```
qsub lidc.sh
```

Launch the evaluation job:

```
# Evaluate will hold until the lidc job is done.
qsub evaluate.sh
```

Use the [elastic load balancer](http://star.mit.edu/cluster/docs/latest/manual/load_balancer.html) to ramp up nodes.

```
# -m 10     -- maximum number of nodes
# -a 2      -- add 2 nodes at a time
# -w 300    -- start adding nodes after jobs have waited 5 minutes (5*60 seconds)
starcluster loadbalance -w 300 -m 10 -a 2 lidc
```

## Stop the nodes

```bash
# Force remove nodes 1-9
starcluster removenode --force lidc node001 node002 node003 node004 node005 node006 node007 node008 node009
```

## Get data from the nodes

```
starcluster get --user sgeadmin lidc lidc.db .
```

### Fixing problems

