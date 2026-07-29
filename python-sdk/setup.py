from setuptools import setup, find_packages

setup(
    name="stuntdouble",
    version="0.1.0",
    description="Python wrapper around the StuntDouble container sandbox CLI",
    author="StuntDouble Team",
    packages=find_packages(),
    install_requires=[],
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
)
