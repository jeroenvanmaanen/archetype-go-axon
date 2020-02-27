#!/usr/bin/false

DOCKER_REPOSITORY='jeroenvm'
EXAMPLE_IMAGE_VERSION='0.0.1-SNAPSHOT'
UI_SERVER_PORT='3000'
API_SERVER_PORT='8181'
AXON_SERVER_PORT='8024'
AXON_VERSION='4.0'

EXTRA_VOLUMES="
      -
        type: bind
        source: ${PROJECT}
        target: ${PROJECT}"
