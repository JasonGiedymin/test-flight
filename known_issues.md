Known Issues
------------

## Scenario: As a user

1. Error while trying to communicate to docker endpoint: http://some.ip:someport:

  This can occur if no docker command was supplied in the Dockerfile, though
  remember that the Dockerfile is created via a template. Did you muck with it?

## Scenario: As a developer

1. Modifications to local code not being reflected when running `make` commands

  `make install` will install compiled code into `$GOPATH/pkg`. When running
  locally looks like this compiled code (in particular `lib`) is picked up
  first. 

  The solution is to delete all references of test-flight from `$GOPATH/pkg`.
  The full path is: 
  
    `/pkg/<your-arch-and-os>/github.com/JasonGiedymin/test-flight`

  You can try to delete it via:

    rm -R `/pkg/<your-arch-and-os>/github.com/JasonGiedymin/test-flight`
