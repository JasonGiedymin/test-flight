Notes
-----

## Enterprise Features

1. Multiple Languages in buildfile. Allows choosing of OracleJDK vs OpenJDK or 
   different c++ compilers.
1. Registry account allowing custom search and docker reference instead of
   public one.
1. Hosted registry account on premise.

## Plan B

Now: Travis -> BuildFile -> Ansible -> Docker
converts travis files to voom build files and runs everything in docker.

Future?: Travis -> SimpleBuildFile -> Docker
converts travis files to a simple buildfile which executes in docker

Now:
  - can test ansible scripts
  - can run ansible on any os, including dockers, has more expressive scripting
  - can test software deployment
  - can test source code compilation
  - can produce a docker artifact at the end*

Future only uses docker to test code. Tests don't test deployments, only
compilation.
  - can only handle compilation, script must work for all operating systems.
    *this is the failure mode, where in multiple OS's you cannot use the same
    script. This could also apply to different base docker images where each
    image is based on different distros.
