# Scaleway Function for Kapsule and Database interactions

This is a Scaleway Function in Golang that sets automatically your Scaleway database's security group in order to allow your Kapsule nodes' IP addresses.

It is important to note that this function as it is, will overwrite your existing security groups.

## How to use it
1. Clone the repository:
2. Compress the content:

`cd kapsule-database-interactions`

`zip archive.zip vendor main.go go.mod go.sum`

3. Create a Scaleway function with the compressed file using the environment variables: ([see the documentation](https://www.scaleway.com/en/docs/compute/functions/quickstart/) for further help) 

ORGANIZATION_ID: your organization ID
ACCESS_KEY: your API access key
SECRET_KEY: your API secret key
DATABASE_INSTANCE_ID: your Scaleway Database's instance ID
KAPSULE_CLUSTER_ID: your Scaleway Kapsule's cluster ID
REGION: the region of your infrastructure

4. Call the function with its endpoint

                         

