# Docker ELK stack

[![Join the chat at https://gitter.im/deviantony/docker-elk](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/deviantony/docker-elk?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Run the latest version of the ELK (Elasticsearch, Logstash, Kibana) stack with Docker and Docker-compose.

It will give you the ability to analyze any data set by using the searching/aggregation capabilities of Elasticsearch and the visualization power of Kibana.

Based on the official images:

* [elasticsearch](https://github.com/elastic/elasticsearch-docker)
* [logstash](https://github.com/elastic/logstash-docker)
* [kibana](https://github.com/elastic/kibana-docker)

**Note**: Other branches in this project are available:

* ELK 5 with X-Pack support: https://github.com/deviantony/docker-elk/tree/x-pack
* ELK 5 in Vagrant: https://github.com/deviantony/docker-elk/tree/vagrant
* ELK 5 with Search Guard: https://github.com/deviantony/docker-elk/tree/searchguard

# Requirements

## Setup

1. Install [Docker](http://docker.io).
2. Install [Docker-compose](http://docs.docker.com/compose/install/) **version >= 1.6**.
3. Clone this repository

## Increase `vm.max_map_count` on your host

You need to increase the `vm.max_map_count` kernel setting on your Docker host.
To do this follow the recommended instructions from the Elastic documentation: [Install Elasticsearch with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode)

## SELinux

On distributions which have SELinux enabled out-of-the-box you will need to either re-context the files or set SELinux into Permissive mode in order for docker-elk to start properly.
For example on Redhat and CentOS, the following will apply the proper context:

```bash
$ chcon -R system_u:object_r:admin_home_t:s0 docker-elk/
```

# Usage

Start the ELK stack using *docker-compose*:

```bash
$ docker-compose up
```

You can also choose to run it in background (detached mode):

```bash
$ docker-compose up -d
```

Now that the stack is running, you'll want to inject logs in it. The shipped logstash configuration allows you to send content via tcp:

```bash
$ nc localhost 5000 < /path/to/logfile.log
```

And then access Kibana UI by hitting [http://localhost:5601](http://localhost:5601) with a web browser.

*NOTE*: You'll need to inject data into logstash before being able to configure a logstash index pattern in Kibana. Then all you should have to do is to hit the create button.

Refer to [Connect Kibana with Elasticsearch](https://www.elastic.co/guide/en/kibana/current/connect-to-elasticsearch.html) for detailed instructions about the index pattern configuration.

By default, the stack exposes the following ports:
* 5000: Logstash TCP input.
* 9200: Elasticsearch HTTP
* 9300: Elasticsearch TCP transport
* 5601: Kibana

*WARNING*: If you're using *boot2docker*, you must access it via the *boot2docker* IP address instead of *localhost*.

*WARNING*: If you're using *Docker Toolbox*, you must access it via the *docker-machine* IP address instead of *localhost*.

# Configuration

*NOTE*: Configuration is not dynamically reloaded, you will need to restart the stack after any change in the configuration of a component.

## How can I tune Kibana configuration?

The Kibana default configuration is stored in `kibana/config/kibana.yml`.

## How can I tune Logstash configuration?

The logstash configuration is stored in `logstash/config/logstash.yml`.

It is also possible to map the entire `config` directory inside the container in the `docker-compose.yml`. Update the logstash container declaration to:

```yml
logstash:
  build: logstash/
  volumes:
    - ./logstash/pipeline:/usr/share/logstash/pipeline
    - ./logstash/config:/usr/share/logstash/config
  ports:
    - "5000:5000"
  networks:
    - docker_elk
  depends_on:
    - elasticsearch
```

In the above example the folder `logstash/config` is mapped onto the container `/usr/share/logstash/config` so you can create more than one file in that folder if you'd like to. However, you must be aware that config files will be read from the directory in alphabetical order, and that Logstash will be expecting a [`log4j2.properties`](https://github.com/elastic/logstash-docker/tree/master/build/logstash/config) file for its own logging.

## How can I specify the amount of memory used by Logstash?

The Logstash container use the *LS_HEAP_SIZE* environment variable to determine how much memory should be associated to the JVM heap memory (defaults to 500m).

If you want to override the default configuration, add the *LS_HEAP_SIZE* environment variable to the container in the `docker-compose.yml`:

```yml
logstash:
  build: logstash/
  volumes:
    - ./logstash/pipeline:/usr/share/logstash/pipeline
  ports:
    - "5000:5000"
  networks:
    - docker_elk
  depends_on:
    - elasticsearch
  environment:
    - LS_HEAP_SIZE=2048m
```

## How can I add Logstash plugins? ##

To add plugins to logstash you have to:

1. Add a RUN statement to the `logstash/Dockerfile` (ex. `RUN logstash-plugin install logstash-filter-json`)
2. Add the associated plugin code configuration to the `logstash/pipeline/logstash.conf` file

## How can I enable a remote JMX connection to Logstash?

As for the Java heap memory, another environment variable allows to specify JAVA_OPTS used by Logstash. You'll need to specify the appropriate options to enable JMX and map the JMX port on the docker host.

Update the container in the `docker-compose.yml` to add the *LS_JAVA_OPTS* environment variable with the following content (I've mapped the JMX service on the port 18080, you can change that), do not forget to update the *-Djava.rmi.server.hostname* option with the IP address of your Docker host (replace **DOCKER_HOST_IP**):

```yml
logstash:
  build: logstash/
  volumes:
    - ./logstash/pipeline:/usr/share/logstash/pipeline
  ports:
    - "5000:5000"
  networks:
    - docker_elk
  depends_on:
    - elasticsearch
  environment:
    - LS_JAVA_OPTS=-Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.port=18080 -Dcom.sun.management.jmxremote.rmi.port=18080 -Djava.rmi.server.hostname=DOCKER_HOST_IP -Dcom.sun.management.jmxremote.local.only=false
```

## How can I tune Elasticsearch configuration?

The Elasticsearch container is using the [shipped configuration](https://github.com/elastic/elasticsearch-docker/blob/master/build/elasticsearch/elasticsearch.yml).

If you want to override the default configuration, create a file `elasticsearch/config/elasticsearch.yml` and add your configuration in it.

Then, you'll need to map your configuration file inside the container in the `docker-compose.yml`. Update the elasticsearch container declaration to:

```yml
elasticsearch:
  build: elasticsearch/
  ports:
    - "9200:9200"
    - "9300:9300"
  environment:
    ES_JAVA_OPTS: "-Xms1g -Xmx1g"
  networks:
    - docker_elk
  volumes:
    - ./elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
```

You can also specify the options you want to override directly via environment variables:

```yml
elasticsearch:
  build: elasticsearch/
  ports:
    - "9200:9200"
    - "9300:9300"
  environment:
    ES_JAVA_OPTS: "-Xms1g -Xmx1g"
    network.host: "_non_loopback_"
    cluster.name: "my-cluster"
  networks:
    - docker_elk
```

## How can I scale up the Elasticsearch cluster?

Follow the instructions from the Wiki: [Scaling up Elasticsearch](https://github.com/deviantony/docker-elk/wiki/Elasticsearch-cluster)

# Storage

## How can I store Elasticsearch data?

The data stored in Elasticsearch will be persisted after container reboot but not after container removal.

In order to persist Elasticsearch data even after removing the Elasticsearch container, you'll have to mount a volume on your Docker host. Update the elasticsearch container declaration to:

```yml
elasticsearch:
  build: elasticsearch/
  ports:
    - "9200:9200"
    - "9300:9300"
  environment:
    ES_JAVA_OPTS: "-Xms1g -Xmx1g"
    network.host: "_non_loopback_"
    cluster.name: "my-cluster"
  networks:
    - docker_elk
  volumes:
    - /path/to/storage:/usr/share/elasticsearch/data
```

This will store elasticsearch data inside `/path/to/storage`.
# NGINX-Modular
