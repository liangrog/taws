taws
====
DevOps tool for interacting with AWS resources

Ansible
-------
`get-inventory`: Get EC2 instance inventory for ansible:

    Usage:
      taws ansible get-inventory [flags]

    Flags:
          --filter-by string      Filter result by Available filters: tags, asg-name
          --filter-value string   Filter values. if filtered by tags, use 'key1=value1;key2=value2' format. If filtered by asg-name, use 'name' string
          --group-by string       Group result by. Available: asg
      -h, --help                  help for get-inventory
          --to-file string        Full file path for alternative inventory file (default to "./ec2-inventory")
          --use-public-ip         If to use public IP rather than private IP

    Global Flags:
          --profile string   AWS CLI profile name
          --region string    AWS region to access


To user tag filter:

    $ taws ansible get-inventory --filter-by tags --filter-value key1=value1;key2=value 

To be grouped by autoscaling group:

    $ taws ansible get-inventory --group-by asg

To be grouped by autoscaling group and filtered by group name

    $ taws ansible get-inventory --group-by asg --filter-by asg-name --filter-value XXX

To only output public IP address (note: instances have no public IP will be ignored):

    $ taws ansible get-inventory --use-public-ip
