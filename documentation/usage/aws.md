# Processing on AWS

Processing the required data (1300+ studies) can be efficiently done using cluster computing.  A [StarCluster](http://star.mit.edu/cluster/docs/latest/index.html) is an on-demand AWS cluster of any size.

## Pricing analysis

The LIDC data is approximately 150G.  EBS storage is $0.05 per GB-month.  To store the LIDC data would cost $7.50 / month on Magnetic volumes, $15.00 / month using SSD.  Processing will likely take 3x, so we'll allocate 500G at a cost of $25.00 / month.

Spot instances are a great value, a `c4.large` has 2 vCPU and 3.75G memory and typically runs ~ $0.01 / hr.  Requesting spot instances is easy, set `spot_bid = 0.02` (or the maximum bid).

### Free tier

- Gives 750 hrs of `t2.micro` or `t1.micro`
- Gives 30 GB EBS storage / month (SSD or magnetic)

## Install StarCluster locally

[StarCluster](http://star.mit.edu/cluster/docs/latest/index.html) is a project to make mananging a OpenGridEngine cluster on AWS easy.

```bash
# Install Python's virtualenv support
pip install --user virtualenv

# Create the virtualenv in the local directory
virtualenv venv

# Activate the local virtualenv
source venv/bin/activate

# Install StarCluster in the virtualenv
easy_install StarCluster

# Install plotting software, didn't work
# easy_install matplotlib
# easy_install numpy
```

Cool! Now StarCluster is installed and we can do interesting things with it.

## Configure

Setup the configuration:

```bash
starcluster help
```

Select `2` to write the config file.  The edit according to the [Quick Start guide](http://star.mit.edu/cluster/docs/latest/quickstart.html).  Using a 2 node cluster of `t2.micro` instances to work with the AMI instance and conform to the AWS free tier.

Using a non-root account called `cluster`.  Created a group called `starcluster` and gave it EC2 and IAM access.

### Specific configuration changes

```
[aws info]
AWS_ACCESS_KEY_ID = ************* #your_aws_access_key_id
AWS_SECRET_ACCESS_KEY =  ***************** #your_secret_access_key
# replace this with your account number
AWS_USER_ID= ********* #your userid

[cluster smallcluster]
CLUSTER_SIZE = 2
NODE_IMAGE = ami-3393a45a
# Instance type, change later
NODE_INSTANCE_TYPE = t1.micro

# Use the volume
VOLUMES = data,software

# Create an EBS volumes
[volume data]
VOLUME_ID = vol-*****
MOUNT_PATH = /home

[volume software]
VOLUME_ID = vol-*****
MOUNT_PATH = /software
```

### Create a volume

[Creating and formatting](http://star.mit.edu/cluster/docs/latest/manual/volumes.html) an EBS volume is relatively easy:

```bash
starcluster createvolume --name=lidc-data --shutdown-volume-host 300 us-east-1a
starcluster createvolume --name=lidc-software --shutdown-volume-host 8 us-east-1a
```

Creates a `10 GB` volume named `lidc-data` in the `us-east-1a` zone, shutting down the creation host afterward.  This command may also bid on a spot instance for $0.10 with the `--bid 0.10` flag.  The bid is not necessary for a `t1.micro` instance, because it cost $0.05 / hr.

### Resize a volume

[Resizing a volume](http://star.mit.edu/cluster/docs/latest/manual/volumes.html#resizing-volumes) is very easy.

```bash
# Gather some information
starcluster listvolumes
# Do the resze
starcluster resizevolume --shutdown-volume-host vol-##### 20
```

This resizes the volume `vol-#####` to 20 Gb and will shutdown the volume host when the command completes.

starcluster resizevolume --shutdown-volume-host vol-7340df8e 250


### Create a keypair:

```
starcluster createkey mykey -o ~/.ssh/mykey.rsa
```

And started the cluster:

```
starcluster start test

# Start a bigger cluster, 8 nodes 4 CPU / 7.5G 
# starcluster start -c lidc lidc

```

Log in:

```
# As root
starcluster sshmaster test

# As sgeadmin
starcluster sshmaster -u sgeadmin test
```

## Do a little test

```
starcluster sshmaster -u sgeadmin test
cat > sleep.pbs <<EOF
#!/bin/sh
 
for i in {1..60} ; do
       echo $i
       sleep 5
done
EOF

chmod 755 sleep.pbs

# submit
for i in {1..5} ; do
  qsub -o sleep.\$JOB_ID.log -j yes sleep.pbs
done

# watch
watch qstat -f

```

## Shutdown the cluster

If the cluster is EBS backed, it can be safely shutdown and restarted with all disks stored on EBS.

```
starcluster stop test
```

To fully delete the cluster, terminate it:

```
# Poof!
starcluster terminate test
```

## Vagrant / Building Tooling for Linux

Create a Vagrant box to build LIDCTooling.

```
vagrant init ubuntu/precise64; vagrant up --provider virtualbox
```

Build instructions are found in `buildVagrant.sh`, and result in software installed in `ClusterSoftware`.

## Copy software to StarCluster

```
starcluster put test --user sgeadmin ClusterSoftware /software/ClusterSoftware
```


## Copy results to an S3 container

Using the [AWS command line tools](https://aws.amazon.com/cli/), start up the cluster and install on the head node.  **Important:** First, the `cluster` user defined on AWS must have Read/Write S3 access.  Visit the [IAM page](https://console.aws.amazon.com/iam/home?region=us-east-1#users/cluster) of the [AWS console](https://console.aws.amazon.com) to grant.  Click on `Attach Policy` and choose `AmazonS3FullAccess`.

```bash
# Spin it up
starcluster start lidc
# Install on head node
starcluster sshmaster lidc
pip install awscli

# Back as sgeadmin
starcluster sshmaster --user sgeadmin lidc

# Configure AWS CLI
$ sgeadmin@master:~$ aws configure
AWS Access Key ID [None]: XXXXXXXXXXXXX
AWS Secret Access Key [None]: YYYYYYYYYYYYYYYYYYYYYYYYY
Default region name [None]: <Return>
Default output format [None]: <Return>
```

**Create the bucket and copy**
```bash
aws s3 mb s3://lidc

# Sync segmented
aws s3 sync segmented s3://lidc/segmented

# Create .tgz of results
tar fcvz lidc.tgz lidc.db segmented/
```


