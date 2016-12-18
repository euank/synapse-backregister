# Synapse backregister

Synapse backregister provides a password reset interface to a synapse server via the the shared registration secret.

It allows hosting an endpoint behind some other form of security, such as
client-cert auth or basic auth, that a user can use to register an account on
the given homeserver.

## Usage

A partial example of running it on Kubernetes is available [here](https://github.com/euank/ek8s/tree/e779a5ae2d30b0d9c61cfcd6b45f5e6a5b13129a/matrix/registration).

Of note, the `SYNAPSE_SERVER` environment variable should be set to the
location your homeserver is serving requests on (e.g.
`https://mydomain.tld:8448`), and the `SYNAPSE_SECRET` environment variable
should be set to the value of the `registration_shared_secret` key from the
synapse homeserver's configuration file.

The server will listen on :8080 for requests.

It is recommended that the docker container, available at
`euank/synapse-backregister:latest` be used.

## Contributing

This project was created as a quick and hacky solution to a problem I had. As
such, it's not very flexible and the code isn't very pretty, modular, etc.

Changes to make it more generic or more robust (e.g. moving to a json api from
posting forms, breaking out modules, using a better framework, etc etc) would
be acceptable.

As is usual, opening an issue to discuss any broad or complex changes prior to
starting them is a great way to make sure time isn't wasted on something that
wouldn't fit in.

## License

Apache 2.0
