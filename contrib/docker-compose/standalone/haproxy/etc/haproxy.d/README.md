# Sample haproxy configuration for keyservers
# (c) Martin Dobrev, Andrew Gallagher 2023

# The files in this directory are optimised for hockeypuck's `docker-compose/standalone` deployment.

# You MUST rename the two `*LOCAL*.cfg.tmpl` files by copying them to the corresponding `*.cfg`.
# Hockeypuck's `docker-compose/standalone/mkconfig.bash` script will do this automatically.
# (see `docker-compose/standalone/README.md` for full instructions).
# This allows you to edit the `*LOCAL*.cfg` files without introducing git merge conflicts.

# To facilitate portability, these files have been parameterised using envar substitution, i.e. "${...}"
# You can populate these values at runtime by setting the corresponding environment variables:

# FQDN                  the FQDN of this server, note that aliases must also be configured
# CERTBOT_HOST_PORT     backend for ACME requests, in the form `host:port`
# PROMETHEUS_HOST_PORT  backend for prometheus monitoring, in the form `host:port`
# KEYSERVER_HOST_PORT   backend for the keyserver, in the form `host:port`
# HAP_CONF_DIR          location of config files
#                       normally `/etc/haproxy` for baremetal, `/usr/local/etc/haproxy` for docker
#                       it must have a subdir `lists` containing `blacklist.list` and `whitelist.list` (can be empty files)
# HAP_CACHE_DIR         persistent state store, must contain `tor_exit_relays.list` (refreshed externally)
# HAP_CERT_DIR          parent directory of SSL/TLS certificate directory
#                       it must contain a subdirectory named after the FQDN, itself containing `fullchain.pem` and `fullchain.pem.key`
#                       e.g. for letsencrypt this will be `/etc/letsencrypt/live`
# HAP_DHPARAM_FILE      Diffie-Hellman parameters for SSL/TLS

# These envars are normally supplied by `docker-compose/standalone/docker-compose.yml` from `docker-compose/standalone/.env`.

# Note that after the cache files or the SSL certs are updated externally, haproxy should be soft reloaded by sending it a HUP signal.
