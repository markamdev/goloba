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

TODO

## Testing

TODO

## Licensing

Code is published under [MIT License](https://opensource.org/licenses/MIT) as it seems to be the most permissive license. If for some reason you need to have this code published with other license (ex. to reuse the code in your project) please contact [author](#author-/-contact) directly.

## Author / Contact

If you need to contact me feel free to write me an email:  
[markamdev.84#dontwantSPAM#gmail.com](mailto:)
