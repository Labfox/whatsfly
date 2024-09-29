import logging

import requests
from setuptools.command.install import install
import subprocess
import platform
import os
from .github_actions_download import download_file


def get_dll_filename(h=False, version=None):
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

    if version != None:
        if not h:
            return f"whatsmeow/whatsmeow-{current_os}-{go_arch}-{version}.{dll_extension}"
        else:
            return f"whatsmeow/whatsmeow-{current_os}-{go_arch}-{version}.h"

    if not h:
        return f"whatsmeow/whatsmeow-{current_os}-{go_arch}.{dll_extension}"
    else:
        return f"whatsmeow/whatsmeow-{current_os}-{go_arch}.h"


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
    env = os.environ.copy()
    env["GOOS"] = current_os
    env["GOARCH"] = go_arch

    go_build_cmd = [
        "go",
        "build",
        "-buildmode=c-shared",
        "-o",
        f"whatsmeow/whatsmeow-{current_os}-{go_arch}.{dll_extension}",
        "main.go",
    ]
    logging.debug(
        f"building Go module with command: {' '.join(go_build_cmd)} in directory {os.getcwd()}/whatsfly/dependencies"
    )

    root_dir = os.path.abspath(os.path.dirname(__file__))

    # Run the Go build command
    status_code = subprocess.check_call(go_build_cmd, cwd=f"whatsfly/dependencies")
    logging.debug(f"Go build command exited with status code: {status_code}")
    if status_code == 127:
        raise RuntimeError("Go build impossible")
    if status_code != 0:
        raise RuntimeError("Go build failed - this package cannot be installed")


def ensureUsableBinaries():
    root_dir = os.path.abspath(os.path.dirname(__file__))

    try:
        import whatsfly.whatsmeow
        return
    except OSError:
        logging.info("Binary unexisent, trying to build")

    os.mkdir(root_dir+"/whatsmeow")

    try:
        build()
        import whatsfly.whatsmeow
        return
    except FileNotFoundError:
        logging.info("Go unusable")
    except RuntimeError:
        logging.warning("Unexpected error while building")

    logging.info("Trying to download pre-built binaries")

    download_file(
        get_dll_filename(version="v20").replace("whatsfly/", "").replace("whatsmeow/", ""),
        root_dir.replace("dependencies", "")+"/dependencies/whatsmeow/"+get_dll_filename().replace("whatsfly/", "").replace("whatsmeow/", ""),
        version="v20"
    )

    download_file(
        get_dll_filename(h=True, version="v20").replace("whatsfly/", "").replace("whatsmeow/", ""),
        root_dir.replace("dependencies", "") + "/dependencies/whatsmeow/" + get_dll_filename(h=True).replace("whatsfly/",
                                                                                                       "").replace(
            "whatsmeow/", ""),
        version="v20"
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
