#!/usr/bin/env python
from setuptools import setup, find_packages
from codecs import open
from setuptools.command.install import install
import subprocess
import platform
import os

setup(
    name="whatsfly",
    version="0.1.3",
    license="MIT",
    author="Doy Bachtiar, Otamay, David Arnold, LabFox, Ivo Bellin Salarin",
    author_email="adityabachtiar996@gmail.com, mauricio@ulisse.io,  labfoxdev@gmail.com, ivo@nilleb.com",
    url="https://whatsfly.labfox.fr",
    keywords="whatsfly whatsapp python",
    description="WhatsApp on the fly.",
    packages=find_packages(),
    install_reqs = ["types-PyYAML", "setuptools", "requests", "qrcode"],
    include_package_data=True,
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
