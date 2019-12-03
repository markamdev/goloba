# GoLoBa - simple Go(lang) Load Balancer

This project is an implementation of a simple load balancer and dummy HTTP server application (for demo/testing purposes).

# Introduction

Application is written totaly in Go (Google Golang) with use of standard library in any possible place. It's usability is limited to load balancing demonstration in easiest possible way. The production use of this project is not recommended but feel free to use the code and final binaries in any presentation or academic project you want.

# Licensing
Code is published under [MIT License](https://opensource.org/licenses/MIT) as it seems to be the most permissive license. If for some reason you need to have this code published with other license (to be honest: I can't imagine why) please contact author directly.

# Building

To build GoLoBa (together with dummyserver test application) just call *make* in main directory. It will create *build* subdir containing output binaries and sample configuration file:

    $ make
    BUILD GoLoBa
    COPY config
    BUILD dummyserver

# Usage

**TODO**

# Testing

**TODO**

# Author / Contact

If you need to contact me feel free to write me an email:  
[markamdev.84#dontwantSPAM#gmail.com]()