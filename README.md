taws
====
DevOps tool for interacting with AWS resources

Ansible
-------
Get EC2 instance inventory by lookup all autoscaling groups: 

    $ taws ansible get-inventory

This will return all EC2 instances in autoscaling groups. To filter it down, you can use filter on autoscaling group name:

    $ taws ansible get-inventory --filter-name xxxx

To use alternative output file, run:

    $ taws ansible get-inventory --inventory-file /home/xxx/yyy
