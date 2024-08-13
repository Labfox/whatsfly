import requests
from setuptools.command.install import install
import subprocess
import platform
import os


def get_dll_filename(h=False):
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
    print(
        f"building Go module with command: {' '.join(go_build_cmd)} in directory {os.getcwd()}/whatsfly/dependencies"
    )

    # Run the Go build command
    status_code = subprocess.check_call(go_build_cmd, cwd="whatsfly/dependencies")
    print(f"Go build command exited with status code: {status_code}")
    if status_code == 127:
        raise RuntimeError("Go build impossible")
    if status_code != 0:
        raise RuntimeError("Go build failed - this package cannot be installed")


def ensureUsableBinaries():
    try:
        import whatsfly.whatsmeow
        return
    except OSError:
        print("Binary unexisent, trying to build")

    try:
        build()
        import whatsfly.whatsmeow
        return
    except FileNotFoundError:
        print("Go unusable")
    except RuntimeError:
        print("Error while building")

    print("Trying to download pre-built binaries")
    url = f"https://github.com/Labfox/whatsfly/raw/jit-compilation/whatsfly/dependencies/{get_dll_filename().replace("whatsmeow/", "whatsmeow/static/")}"
    h_url = f"https://github.com/Labfox/whatsfly/raw/jit-compilation/whatsfly/dependencies/{get_dll_filename(h=True).replace("whatsmeow/", "whatsmeow/static/")}"

    print(f"Dowloading {url} and {h_url}")

    try:
        rq = requests.get(url, stream=True)
        if rq.status_code != 200:
            raise RuntimeError(
                f"Server responded with {rq.status_code}, impossible to find the binaries, giving up"
            )
        open(f"whatsfly/dependencies/{get_dll_filename()}", "wb").write(rq.content)

        rq = requests.get(url, stream=True)
        if rq.status_code != 200:
            raise RuntimeError(
                f"Server responded with {rq.status_code}, impossible to find the binaries, giving up"
            )
        open(f"whatsfly/dependencies/{get_dll_filename(h=True)}", "wb").write(
            rq.content
        )
    except Exception:
        raise RuntimeError("Impossible to find the binaries, giving up")


class BuildGoModule(install):
    def run(self):
        # Ensure the Go module is built before the Python package
        self.build_go_module()
        super().run()

    def build_go_module(self):
        try:
            build()
        except RuntimeError:
            print("Build unsuccessful, will retry on runtime")
