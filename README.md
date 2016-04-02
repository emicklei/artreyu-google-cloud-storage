artreyu-google-cloud-storage
===
plugin for artreyu

### Requirements

- artreyu
- gcloud SDK
- bucket on Google Cloud Storage

### Build

	make local

### Usage

Inside a folder that contains an `artreyu.yaml` file.

	artreyu archive -r gcs someasset.zip

#### Configuration

To use this plugin, add a repository to your configuration file `.artreyu`

	repositories:	
	- name:	    gcs
	  plugin:   gcs
	  bucket:	yours_assets
	  path:     /
	
The plugin is called `gcs`. You need to create a bucket first.