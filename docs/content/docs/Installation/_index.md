+++
date = '2024-12-24T21:25:42+01:00'
draft = false
title = 'Installation'
weight = 10
+++
# Installation Guide for Wait4it

To get started with Wait4it, you can either download the latest release or build it yourself. Follow the instructions below based on your preference.

## Option 1: Download the Latest Release

The easiest way to get started is by downloading the latest release of Wait4it.

1. Visit the [Wait4it Releases Page](https://github.com/ph4r5h4d/wait4it/releases).
2. Download the appropriate binary for your operating system (Linux, macOS, or Windows).
3. Extract the downloaded file and move it to a directory of your choice.
4. Ensure that the directory is included in your system's `PATH` to run Wait4it from anywhere in your terminal.

## Option 2: Build Wait4it Yourself

If you prefer to build Wait4it from source, follow these steps:

### Prerequisites

Make sure you have [Go](https://golang.org/dl/) installed on your machine.

### Build Instructions

1. Clone the Wait4it repository using `git` or the GitHub CLI:

    - Using `git`:
      ```bash
      git clone https://github.com/ph4r5h4d/wait4it.git
      ```

    - Using GitHub CLI (`gh`):
      ```bash
      gh repo clone ph4r5h4d/wait4it
      ```

2. Navigate to the project directory:
    ```bash
    cd wait4it
    ```

3. Build the project:
    ```bash
    go build
    ```

4. After the build process completes, the `wait4it` binary will be created in the current directory.

5. Move the binary to a directory in your `PATH` for easy access:
    ```bash
    mv wait4it /usr/local/bin/
    ```

6. Verify the installation by running:
    ```bash
    wait4it --version
    ```

Now you're ready to start using Wait4it!

---

For more details on how to use Wait4it, refer to the left sidebar for specific checks and their usage instructions. 

