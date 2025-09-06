#!/usr/bin/env python
from setuptools import setup

setup(
    name="whatsfly",
    version="0.3.2",
    author="LabFox",
    author_email="labfoxdev@gmail.com",
    url="https://whatsfly.labfox.fr",
    keywords="whatsfly whatsapp python",
    description="WhatsApp on the fly.",
    package_dir={"whatsfly": "whatsfly"},
    packages=["whatsfly"],
    install_requires = ["types-PyYAML", "setuptools", "requests", "qrcode", "protobuf"],
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
        "License :: OSI Approved :: MIT License",
    ],
)
