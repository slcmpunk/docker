### Docker users, see [Moby and Docker](https://mobyproject.org/#moby-and-docker) to clarify the relationship between the projects

## Red Hat Patches
Red Hat is carrying a series of experimental patches that we feel are required
for our customers or for our support engineering.

#### Ignore-invalid-host-header-between-go1.6-and-old-docker-clients.patch

https://bugzilla.redhat.com/show_bug.cgi?id=1324150
https://github.com/docker/docker/pull/22000
https://github.com/docker/docker/issues/20865
https://github.com/docker/docker/pull/21423

#### The-following-syscalls-should-not-be-blocked-by-seccomp.patch

Capabilities block these syscalls.

mount, umount2, unshare, reboot and name\_to\_handle\_at are all needed to
run systemd as pid1 in a container, they work fine with sys\_admin disabled
and have functionality in the kernel that is available to a non privileged
process.  There is no easy way to discover which syscalls are blocked, so
we end up more likely with the user doing a --privileged.

With UserNamespace we want to allow users to potentially setup unshare additional
namespaces.

man reboot
...
Behavior inside PID namespaces
Since Linux 3.4, when reboot() is called from a PID namespace (see
		pid_namespaces(7)) other than the initial PID namespace, the effect
		of the call is to send a signal to the namespace "init" process.
		LINUX_REBOOT_CMD_RESTART and LINUX_REBOOT_CMD_RESTART2 cause a SIGHUP
		signal to be sent.  LINUX_REBOOT_CMD_POWER_OFF and
		LINUX_REBOOT_CMD_HALT cause a SIGINT signal to be sent.

https://github.com/docker/docker/pull/21287

#### Add-dockerhooks-exec-custom-hooks-for-prestart/poststop-containers.patch

With the addition of runc/hooks support we want to add a feature
to allow third parties to run helper programs before a docker container
gets started and just after the container finishes.

For example we want to add a RegisterMachine hook.

For systems that support systemd/RegisterMachine, this hook would register
a machine to the machinectl.  machinectl could then list docker containers
along with other virtulization environments like kvm, and systemd-nspawn
containers. Overtime we would want to implement other machinectl features
to get docker containers better integrated into the system and machinectl.

Another example of a dockerhook might be for people wanting to do better logging
of starting and stopping of containers.  For example have a log agent that
records when a container starts and stops and then sends a message to a
monitoring station.

Dockerhooks reads directory in either /usr/lib/docker/hooks.d or
/usr/libexec/docker/hooks.d to search for hooks, if the directory exists
docker will execute the executables in this directory via runc/libcontainer i
using PreStart and PostStop.  It will also send the config.json file as the
second paramater.

https://github.com/docker/docker/pull/17021

#### Return-rpm-version-of-packages-in-docker-version.patch

Red Hat Support wants to know the version of the rpm package that docker
is running.  This patch allows the distribution to show the rpm version
of the client and server using the `docker info` command.  Docker upstream
was not interested in this patch. This patch only affects the `docker info`
command.

https://github.com/docker/docker/pull/14591

#### rebase-distribution-specific-build.patch

Current docker tests run totally on Ubuntu.  We want to be able to make sure
all tests run on platforms and containers we support.  This patch allows us
to run distribution specific tests. This means that `make test` on RHEL7 will
use rhel7 based images for the test, and `make test` on Fedora will use
fedora based images for the test.  Docker was not crazy about the substitutions
done in this patch, and the changes never show up in the docker client/daemon
that we ship.

https://github.com/docker/docker/pull/15364

#### Add-RHEL-super-secrets-patch.patch

This patch allows us to provide subscription management information from the
host into containers for both `docker run` and `docker build`.  In order to get
access to RHEL7 and RHEL6 content inside of a container via yum you need to
use a subscription.  This patch allows the subscriptions to work inside of the
container.  Docker thought this was too RHEL specific so told us to carry a
patch.

https://github.com/docker/docker/pull/6075

#### Add-add-registry-and-block-registry-options-to-docke.patch

Customers have been asking for us to provide a mechanism to allow additional
registries to be specified in addition to docker.io.  We also want to have
Red Hat Content available from our Registry by default. This patch allows
users to customize the default registries available.  We believe this closely
aligns with the way yum and apt currently work.  This patch also allows
customers to block images from registries.  Some customers do not want software
to be accidentally pulled and run on a machine.  Some customers also have no
access to the internet and want to setup private registries to handle their
content.

https://github.com/docker/docker/pull/11991
https://github.com/docker/docker/pull/10411

#### Improved-searching-experience.patch

Red Hat wants to allow users to search multiple registries as described above.
This patch improves the search experience.

#### System-logging-for-docker-daemon-API-calls.patch
Red Hat wants to log all access to the docker daemon.  This patch records all
access to the docker daemon in syslog/journald.  Currently docker logs access
in its event logs but these logs are not stored permanently and do not record
critical information like loginuid to record which user created/started/stopped
... containers.  We will continue to work with Docker to get this patch
merged. Docker indicates that they want to wait until the authentication patches
get merged.

https://github.com/docker/docker/pull/14446

#### Audit-logging-support-for-daemon-API-calls.patch

This is a follow on patch to the previous syslog patch, to record all auditable
events to the audit log.  We want to get `docker daemon` to some level of common
criteria.  In order to do this administrator activity must be audited.

https://github.com/rhatdan/docker/pull/109

#### Add-volume-support-to-docker-build.patch

This patch adds the ability to add bind mounts at build-time. This will be
helpful in supporting builds with host files and secrets. The --volume|-v flag
is added for this purpose. Trying to define a volume (ala docker run) errors
out. Each bind mounts' mode will be read-only and it will preserve any SELinux
label which was defined via the cli (:[z,Z]). Defining a read-write mode (:rw)
will just print a warning in the build output and the actual mode will be
changed to read-only.

Docker: the container engine [![Release](https://img.shields.io/github/release/docker/docker.svg)](https://github.com/docker/docker/releases/latest)
============================

### Docker maintainers and contributors, see [Transitioning to Moby](#transitioning-to-moby) for more details

The Moby Project
================

![Moby Project logo](docs/static_files/moby-project-logo.png "The Moby Project")

Moby is an open-source project created by Docker to advance the software containerization movement.
It provides a “Lego set” of dozens of components, the framework for assembling them into custom container-based systems, and a place for all container enthusiasts to experiment and exchange ideas.

# Moby

## Overview

At the core of Moby is a framework to assemble specialized container systems.
It provides:

- A library of containerized components for all vital aspects of a container system: OS, container runtime, orchestration, infrastructure management, networking, storage, security, build, image distribution, etc.
- Tools to assemble the components into runnable artifacts for a variety of platforms and architectures: bare metal (both x86 and Arm); executables for Linux, Mac and Windows; VM images for popular cloud and virtualization providers.
- A set of reference assemblies which can be used as-is, modified, or used as inspiration to create your own.

All Moby components are containers, so creating new components is as easy as building a new OCI-compatible container.

## Principles

Moby is an open project guided by strong principles, but modular, flexible and without too strong an opinion on user experience, so it is open to the community to help set its direction.
The guiding principles are:

- Batteries included but swappable: Moby includes enough components to build fully featured container system, but its modular architecture ensures that most of the components can be swapped by different implementations.
- Usable security: Moby will provide secure defaults without compromising usability.
- Container centric: Moby is built with containers, for running containers.

With Moby, you should be able to describe all the components of your distributed application, from the high-level configuration files down to the kernel you would like to use and build and deploy it easily.

Moby uses [containerd](https://github.com/containerd/containerd) as the default container runtime.

## Audience

Moby is recommended for anyone who wants to assemble a container-based system. This includes:

- Hackers who want to customize or patch their Docker build
- System engineers or integrators building a container system
- Infrastructure providers looking to adapt existing container systems to their environment
- Container enthusiasts who want to experiment with the latest container tech
- Open-source developers looking to test their project in a variety of different systems
- Anyone curious about Docker internals and how it’s built

Moby is NOT recommended for:

- Application developers looking for an easy way to run their applications in containers. We recommend Docker CE instead.
- Enterprise IT and development teams looking for a ready-to-use, commercially supported container platform. We recommend Docker EE instead.
- Anyone curious about containers and looking for an easy way to learn. We recommend the docker.com website instead.

# Transitioning to Moby

Docker is transitioning all of its open source collaborations to the Moby project going forward.
During the transition, all open source activity should continue as usual.

We are proposing the following list of changes:

- splitting up the engine into more open components
- removing the docker UI, SDK etc to keep them in the Docker org
- clarifying that the project is not limited to the engine, but to the assembly of all the individual components of the Docker platform
- open-source new tools & components which we currently use to assemble the Docker product, but could benefit the community
- defining an open, community-centric governance inspired by the Fedora project (a very successful example of balancing the needs of the community with the constraints of the primary corporate sponsor)

-----

Legal
=====

*Brought to you courtesy of our legal counsel. For more context,
please see the [NOTICE](https://github.com/moby/moby/blob/master/NOTICE) document in this repo.*

Use and transfer of Moby may be subject to certain restrictions by the
United States and other governments.

It is your responsibility to ensure that your use and/or transfer does not
violate applicable laws.

For more information, please see https://www.bis.doc.gov


Licensing
=========
Moby is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.
