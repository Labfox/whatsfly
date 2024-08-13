#!/usr/bin/env python
from setuptools import setup, find_packages
from codecs import open
from setuptools.command.install import install
import os
import subprocess
import platform


class BuildGoModule(install):
    def run(self):
        # Ensure the Go module is built before the Python package
        self.build_go_module()
        super().run()

    def build_go_module(self):
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
        go_arch = arch_map.get(current_arch, current_arch)

        # Set the environment variables for Go build
        env = os.environ.copy()
        env["GOOS"] = current_os
        env["GOARCH"] = go_arch

        go_build_cmd = [
            "go",
            "build",
            "-buildmode=c-shared",
            "-o",
            f"whatsmeow/whatsmeow-{current_os}-{go_arch}.dylib",
            "main.go",
        ]
        print(
            f"building Go module with command: {' '.join(go_build_cmd)} in directory {os.getcwd()}/whatsfly/dependencies"
        )
        # Run the Go build command
        status_code = subprocess.check_call(go_build_cmd, cwd="whatsfly/dependencies")
        print(f"Go build command exited with status code: {status_code}")
        if status_code != 0:
            raise RuntimeError("Go build failed - this package cannot be installed")


setup(
    name="whatsfly",
    version="0.1.0",
    license="MIT",
    author="Doy Bachtiar, Lab Fox, Ivo Bellin Salarin",
    author_email="adityabachtiar996@gmail.com, labfoxdev@gmail.com, ivo@nilleb.com",
    url="https://github.com/cloned-doy/whatsfly",
    keywords="whatsfly",
    description="WhatsApp on the fly.",
    long_description=open("README.md", encoding="utf-8").read(),
    long_description_content_type="text/markdown",
    packages=find_packages(),
    include_package_data=True,
    cmdclass={
        "install": BuildGoModule,
    },
    classifiers=[
        "Intended Audience :: Developers",
        "Natural Language :: English",
        "Operating System :: Unix",
        "Operating System :: MacOS :: MacOS X",
        "Operating System :: Microsoft :: Windows",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Environment :: Web Environment",
        "Topic :: Communications",
        "Topic :: Communications :: Chat",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
)
