#TODO
* Add span for FileStorage (?) and S3
* Rewrite storage factory to "pick" a storage-provider based on an env-variables
    * Add support for Azure Blob storage
    * Add support for S3
    * Add whatever GCE uses..
    * Add Dapr support
* Add events to be published on Create,Remove,Update using Kafka
    * Write tests for events
* WHEN DISTRIBUTING THIS AS A DOCKER IMAGE, THIRD-PARTY LICENSES ARE REQUIRED TO BE INCLUDED IN THAT IMAGE!!!

## CLEAN-UP
* Can more thing be passed by reference?
* Can we use goroutines to speed things up?