#TODO
* Add logging to handlers
* Handle "bubbled" errors..
* Only enable OTLP exporter given that environment variables are set
* Add formatters
    * Allow us to set formatter based on env-variable
        * For JSON
        * For Dev
            * Include file-name, method-name, linennumber?
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