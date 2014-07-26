

#
# ------------------------------------------------------
#                       Dockerfile
# ------------------------------------------------------
# image:    test-docker-name1
# tag:      latest
# name:     testflight/test-docker-name1
# version:  v0.1.0
# how-to:   docker build -t testflight/test-docker-name1 .
# Requires: jasongiedymin/ansible-nodejs
#
#       Test-Flight v0.9.8.1 Generated Dockerfile
# ------------------------------------------------------

FROM jasongiedymin/ansible-nodejs
MAINTAINER testflight



# Required Test-Flight Working dir
ENV WORKDIR /tmp/build/test-docker-name1


# -----> Env entries

ENV ENV_MODE dev

# <----- end Env





# -----> Add Simple Entries

ADD tests/test-dirmode/example-playbook/tasks $WORKDIR/tasks

ADD tests/test-dirmode/example-playbook/tests $WORKDIR/tests

ADD tests/test-dirmode/example-playbook/vars $WORKDIR/vars

ADD tests/test-dirmode/example-playbook/handlers $WORKDIR/handlers

# <----- end Add System Entries




# Required Test-Flight test-flight playbook and inventory
ADD tests/test-dirmode/example-playbook/.test-flight $WORKDIR/.test-flight
ADD tests/test-dirmode/example-playbook/.test-flight/inventory /etc/ansible/hosts
ADD tests/test-dirmode/example-playbook/.test-flight/playbook.yml $WORKDIR/playbook.yml




# -----> Install Galaxy Dependencies

RUN ansible-galaxy install JasonGiedymin.ark-nodejs

# <----- end Install Galaxy Dependencies


# Execute Test-Flight Playbook
RUN ansible-playbook $WORKDIR/playbook.yml -c local -vvvv




# -----> Expose entries

Expose 3000

Expose 3001

# <----- end Expose entries



CMD sleep 5 && echo 'Hello Build!'



