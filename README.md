# GoLoBa - Go(lang) Load Balancer

Implementation of a simple load balancer for TCP connections.

## Introduction

**goloba** is a simple example of *Layer 4 balancer* working in *proxy mode*. Each connection established with **goloba** process is forwarded to one of configured endpoints. Endpoint selection is made internally based on implemented balancing algorithm.

Load balancer listening port and list of endpoints (address:port pairs) are provided in configuration file. Please check [Usage](#usage) section for more details.

**goloba** is a demonstration of load balancing idea rather than complete, highly efficient solution and so it is not suggested to use it in any production environment.

## Downloading and building

GoLoBa project is created as an Go module with Makefile provided so can be downloaded by *git clone* to any location and compiled by simple *make*:

```bash
git clone github.com/markamdev/goloba
cd goloba
make
```

Output binaries (**goloba** and **dummyserver** test application) and sample configuration file will be placed inside *build* subdir. **goloba** binary is a load balancer while **dummyserver** is a simple application (based on HTTP server idea) prepared for easy testing ([see Testing](#testing) section).

## Usage

### goloba params and usage

### dummyserver params and usage

### Example of goloba config

## Testing

Project provides two bash scripts that simplifies quick load balancer testing: *start_testbench.sh* and *start_curltest.sh*.

*start_testbench.sh* launches requested number of HTTP servers (**dummyserver** application), prepares *goloba.conf* file and launches **goloba** itself. When testbench is ready (there's no error messages and scripts "freezes" on launched goloba instance) user should call *start_curltest.sh* in separate terminal with appropriate parameters. This second script performs multiple HTTP GET requests (*curl* application needed for this) and prints received web page source on screen. As each dummyserver instance is launched with different message string successive GETs should contain different text.

Test scripts are prepared to be launched on one machine only so testcase does not verify proper forwarding to remote machines.

### start_testbench.sh params and usage

This script gets 3 command line params:

* Listening port for balancer
* Listening port of first HTTP server (first dummyserver instance)
* Number of servers to be launched

To launch balancer listening on port 9000 and 4 servers where first is listening on port 9500 (so rest of servers is listening on ports 9501, 9502 and 9503) user should call:

```bash
./scripts/start_testbench.sh 9000 9500 4
```

### start_curltest.sh params and usage

This script gets 2 command line params:

* Listetning port of balancer
* Number of get requests to be executed

To test environment prepared by *start_testbench.sh* in example above with 5 GET requests user should call:

```bash
./scripts/start_curltest.sh 9000 5
```

As there are only 4 servers launched, result of 1st abd 5th GET request should return same content.

## Licensing

Code is published under [MIT License](https://opensource.org/licenses/MIT) as it seems to be the most permissive license. If for some reason you need to have this code published with other license (ex. to reuse the code in your project) please contact [author](#author-/-contact) directly.

## Author / Contact

If you need to contact me feel free to write me an email:  
[markamdev.84#dontwantSPAM#gmail.com](mailto:)
