#TODO
* Add logging to services and handlers
* Add Spans/activities
* Validate spans and activies in Grafana
* Add metrics
* Validate metrics in Grafana
* Rewrite storage factory to "pick" a storage-provider based on an env-variables
    * Add support for Azure Blob storage
    * Add support for S3
    * Add whatever GCE uses..
* Add events to be published on Create,Remove,Update using Kafka
    * Write tests for events