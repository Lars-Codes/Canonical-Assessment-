# Canonical coding assessment: Bootable image 

This creates a shell script that will run in a Linux environment. It creates and runs an AMD64 Linux Filesystem Image using QEMU that prints "Hello World" after successful 
startup. It is a fully bootable filesystem image. 

---

## Table of Contents

-[Prerequisites](#Prerequisites)
- [Usage](#usage)
- [License](#license)  

## Prerequisites 

Because the assignment specifies that the script will be tested on Ubuntu 20.04 LTS or 22.04 LTS, this script assumes that the following packages are installed: 

bash
sudo
wget
cpio
gzip
qemu-system-x86_64
losetup
mkfs.ext4


## Usage 
Clone or copy the script to your working directory.

Make the script executable:
```
chmod +x build_minimal_linux.sh
```
Run the script:
```
./create-bootable-image.sh

```


## License 
MIT License. [License](https://github.com/Lars-Codes/Canonical-Assessment-/blob/master/LICENSE)

