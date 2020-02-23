# Docker

This directory describes how the components of the application can be
deployed together as a completely functional system.

The setup presupposes local configuration files. Sensible defaults can
be copied from provided samples using `../bin/create-local-settings.sh`.
This does not overwrite existing local files.

Running the complete setup is possible using the single command `./docker-compose-up.sh`.

The dashboard of AxonServer is exposed at [AxonDashboard](http://localhost:8124).