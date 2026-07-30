### Create docker images

Install docker. 

Then build the container from the Dockerfile :

```bash
docker build --build-arg PYTHON_VERSION=3.13 -t py3.13-docker .
```

To build several images, one per python version : 

```bash
docker build --build-arg PYTHON_VERSION=3.11 -t py3.11-docker .
docker build --build-arg PYTHON_VERSION=3.12 -t py3.12-docker .
docker build --build-arg PYTHON_VERSION=3.13 -t py3.13-docker .
```

> The Dockerfile can of course be modified to change the requirements, or add anything you would need in the container

### Use kbox to install contenairized kernels

python kernel :

```bash
kbox install py3.13-docker python3 kbox-py3.13
```

or bash kernel (from root of repo):

```bash
kbox install py3.13-docker ./bash kbox-bash
```
