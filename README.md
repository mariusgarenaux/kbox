# KBox

KBox is a lightweight utility written in **Go**, distributed as a **single standalone binary**. It enables you to run Jupyter kernels inside Docker containers while maintaining a seamless connection via your preferred Jupyter frontend (e.g., Jupyter Notebooks, JupyterLab, or VS Code).

By wrapping the kernel execution in a container, KBox provides:
- **Host Independence**: The host machine does not need to have the kernel runtime installed (e.g., you don't need Python or Bash installed on your host to run a Python or Bash kernel). Everything is fully containerized.
- **Isolated Execution**: Run any kernel (Python, Bash, AI agents, etc.) in a fully containerized environment.
- **Simplified Management**: Easily switch between different environment images and kernel versions without polluting your host machine.
- **Automatic Lifecycle**: Containers are automatically started when the kernel is launched and removed when the notebook is closed.

## Getting Started

### 1. Installation
Compile the project to create the `kbox` binary:

```bash
go build -o kbox ./kbox
```

### 2. Installing a Containerized Kernel
You can create a new kernelspec using either a built-in kernel name or a path to an existing kernelspec directory.

#### Using Built-in Kernels
KBox comes with a set of predefined kernels (like `python3` and `bash`). You can list them using:
```bash
./kbox list
```

To install a built-in kernel:
```bash
./kbox install <image_name:tag> <built_in_name> [new_kernel_name]
```
**Example:**
```bash
./kbox install my-python-env:latest python3 kbox-py3.13
```

#### Using a Custom Kernelspec Path
If you have a specific kernel configuration on your disk, you can point KBox to its directory:
```bash
./kbox install <image_name:tag> /path/to/kernelspec [new_kernel_name]
```

*Note: The Docker image must contain the necessary runtime and kernel package to support the kernel. See [python_dockerfile](./python_dockerfile/) and [ai_agent_dockerfile](./ai_agent_dockerfile/) for examples.*

### 3. Verification
Verify that the new kernel is installed:
```bash
jupyter kernelspec list
```

You should now see your new kernel in the list:
```text
kbox-py3.13    /Users/mgg/Library/Jupyter/kernels/kbox-py3.13
```

## Usage

Once installed, simply select the `kbox-` kernel from your Jupyter dropdown. KBox will launch an instance of the specified Docker image and bridge the communication between the notebook and the container.

**Key Features:**
- **Working Directory**: The directory where the KBox kernel is started is automatically mounted into the container, providing easy access to local datasets and files.
- **Host Independence**: KBox manages the bridge between the host and the container, allowing you to use images that may have different OS requirements than your host.

### Quick Test
You can test the installation via the terminal:
```bash
jupyter console --kernel kbox-py3.13
```
