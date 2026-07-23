from setuptools import setup, find_packages

setup(
    name="stuntdouble",
    version="1.0.0",
    description="Python SDK for the StuntDouble eBPF AI Sandbox",
    author="StuntDouble Team",
    packages=find_packages(),
    install_requires=[],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
)
