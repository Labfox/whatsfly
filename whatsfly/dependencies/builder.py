import logging

import requests
from setuptools.command.install import install
import subprocess
import platform
import os


def download_file(file, path, isSymLink=True):
    if isSymLink:
        r = requests.get(
            f"https://raw.githubusercontent.com/Labfox/whatsfly/refs/heads/prebuilts/{file}"
        )
        if r.status_code != 200:
            raise FileNotFoundError()

        file = r.text

    r2 = requests.get(
        f"https://raw.githubusercontent.com/Labfox/whatsfly/refs/heads/prebuilts/{file}"
    )

    open(path, "wb").write(r2.content)


def get_dll_filename(branch="main"):
    current_os = platform.system().lower()
    current_arch = platform.machine().lower()

    # Map the architecture to Go's naming convention
    arch_map = {
        "x86_64": "amd64",
        "arm64": "arm64",
        "aarch64": "arm64",
    }

    extension_map = {"linux": "so", "windows": "dll", "darwin": "dylib"}

    go_arch = arch_map.get(current_arch, current_arch)
    dll_extension = extension_map.get(current_os, current_os)

    return f"{current_os}-{go_arch}-{branch}/latest.{dll_extension}"


def get_extension_name():
    current_os = platform.system().lower()

    extension_map = {"linux": "so", "windows": "dll", "darwin": "dylib"}

    dll_extension = extension_map.get(current_os, current_os)
    return dll_extension


def build():
    # Define the Go build command, something like
    # GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-darwin-amd64.dylib main.go
    # Detect the current OS and architecture
    current_os = platform.system().lower()
    current_arch = platform.machine().lower()

    # Map the architecture to Go's naming convention
    arch_map = {
        "x86_64": "amd64",
        "arm64": "arm64",
        "aarch64": "arm64",
    }

    extension_map = {"linux": "so", "windows": "dll", "darwin": "dylib"}

    go_arch = arch_map.get(current_arch, current_arch)
    dll_extension = extension_map.get(current_os, current_os)

    # Set the environment variables for Go build

    root_dir = os.path.abspath(os.path.dirname(__file__))

    go_build_cmd = [
        "go",
        "build",
        "-buildmode=c-shared",
        "-o",
        f"{root_dir}/latest.{dll_extension}",
        "main.go",
    ]
    logging.debug(
        f"building Go module with command: {' '.join(go_build_cmd)} in directory {os.getcwd()}/whatsfly/dependencies"
    )

    root_dir = os.path.abspath(os.path.dirname(__file__))

    # Run the Go build command
    status_code = subprocess.check_call(go_build_cmd)
    logging.debug(f"Go build command exited with status code: {status_code}")
    if status_code == 127:
        raise RuntimeError("Go build impossible")
    if status_code != 0:
        raise RuntimeError("Go build failed - this package cannot be installed")


def ensureUsableBinaries():
    branch = "main"
    logging.info("Trying to download pre-built binaries")
    root_dir = os.path.abspath(os.path.dirname(__file__))

    download_file(
        get_dll_filename(branch=branch),
        f"{root_dir}/latest.{get_extension_name()}",
        isSymLink=True,
    )


class BuildGoModule(install):
    def run(self):
        # Ensure the Go module is built before the Python package
        self.build_go_module()
        super().run()

    def build_go_module(self):
        try:
            build()
        except RuntimeError:
            logging.warning("Build unsuccessful, will retry on runtime")


ensureUsableBinaries()
