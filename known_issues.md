Known Issues
------------

1. Error while trying to communicate to docker endpoint: http://some.ip:someport:

  This can occur if no docker command was supplied in the Dockerfile, though
  remember that the Dockerfile is created via a template. Did you muck with it?

