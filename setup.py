#!/usr/bin/env python
from setuptools import setup

setup(
    name="whatsfly",
    version="0.2.0",
    license="MIT",
    author="Doy Bachtiar, Otamay, David Arnold, LabFox, Ivo Bellin Salarin",
    author_email="adityabachtiar996@gmail.com, mauricio@ulisse.io,  labfoxdev@gmail.com, ivo@nilleb.com",
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
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Environment :: Web Environment",
        "Topic :: Communications",
        "Topic :: Communications :: Chat",
        "Topic :: Software Development :: Libraries :: Python Modules",
    ],
)
