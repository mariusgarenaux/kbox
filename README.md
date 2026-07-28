# KBox

Kbox is a command line tool that can manage jupyter kernels started inside docker containers.
The goal is to start jupyter kernels inside a given container, and to send them messages from the host machine.

Kbox helps to manage several kernelspec files, as well as their connection to the containers. For example, kbox could provide a command for starting a one shot container ('jetable'); and also commands to start a new kernel from a running container. It creates a connection file once a container is started, to allow to create a new kernel inside the container.