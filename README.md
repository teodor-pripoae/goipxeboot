# goipxeboot

Server which allows booting custom live Linux distributions via iPXE.

This server handles all parts required for iPXE booting the system, apart from the DHCP server.

## Description

It needs a DHCP server which is able to set up next-server and filename options. The filename option should be set to `ipxe.efi` for UEFI systems and this project currently does not support BIOS systems.

It will spawn a TFTP server on port `69/udp` and a HTTP server on port `8080/tcp` which will serve the iPXE efi executable and the iPXE script. It also supports serving the kernel, initrd and squashfs files via HTTP.

You don't need to have iPXE compatible network card since the TFTP server will serve the iPXE EFI file when booted over traditional PXE.

## Building

To build the project, you need to have a working Go environment. You can build the project by running:

```
$ go build -o goipxeboot ./cmd/goipxeboot
```

[Bazel](https://bazel.build/) is also supported. You can build the project by running:

```
$ bazel build //...
```

## Running

To run the project, you will need to have a DHCP server which is able to set up next-server and filename options. When the server is started it will output the ISC DHCP configuration which you can use to configure your DHCP server.

Example configuration:

```
next-server 10.20.3.4;
filename "ipxe.efi";

class "pxeclients" {
	if exists user-class and option user-class = "iPXE" {
	    filename "http://10.20.3.4:8080/ipxe";
	} else {
	    filename "ipxe.efi";
	}
}
```

You will also need to setup a config file which will allow serving different live Linux systems based on allowed IP addresses. Example configuration:

```yaml
rootDir: "/path/to/root/dir" # Root directory where the filesystems are stored
http:
    ip: "192.168.1.100" # IP of the http server (needed for iPXE script)
    port: 8080 # Port of the http server (needed for iPXE script)
ipxe:
    - name: "example" # Name of the iPXE configuration
      ips: ["192.168.1.124"] # Allowed IPs
      kernelArgs: # Custom kernel arguments used for booting the system
          - "network-config=disabled"
```

You can run the server by running:

```
$ ./goipxeboot server -c goipxeboot.yaml
```

### Filesystem structure

Root directory should contain directories for each iPXE configuration. Each directory should contain the kernel, initrd and squashfs files. Example structure:

```
ipxe
ipxe/ipxe.efi
linux
linux/example
linux/example/vmlinuz
linux/example/squashfs
linux/example/initrd
```

[ipxe.efi](https://boot.ipxe.org/ipxe.efi) can be downloaded from the iPXE website.
