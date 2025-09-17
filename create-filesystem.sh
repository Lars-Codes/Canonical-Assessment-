#!/bin/bash
set -e 

WORKDIR="$(pwd)/workdir"
MNT="$WORKDIR/mnt"
mkdir -p "$WORKDIR" "$MNT"


# Create a virtual disk image 
qemu-img create -f raw hello_canonical.img 1G 

# Attach virtual disk image using losetup 
sudo losetup -fP hello_canonical.img  

# Assign corresponding loop device to variable
device=$(losetup -a 2>/dev/null | grep 'hello_canonical.img' | awk -F: '{print $1}')

# Partition the device 
echo ",,83,*" | sudo sfdisk "$device"

# Refresh partition table 
sudo partprobe "$device"

# Detect first partition when it becomes available 
while ! ls ${device}?* 1> /dev/null 2>&1; do
    sleep 0.1
done 

# Select first parition and assign it to variable 
first_partition=$(ls ${device}?* | head -n1)

# Format the partition 
sudo mkfs.ext4 $first_partition

# Mount the partition locally 
sudo mount "$first_partition" "$MNT"

# Download/install Alpine rootfs tarball for mini root filesystem 
ROOTFS_TAR="$WORKDIR/alpine-minirootfs.tar.gz"
curl -L -o "$ROOTFS_TAR" https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/x86_64/alpine-minirootfs-3.22.1-x86_64.tar.gz

# Extract
sudo tar -xpf "$ROOTFS_TAR" -C "$MNT"

# Mount filesystems 
mount -t proc proc /proc
mount -t sysfs sys /sys
mount -t devtmpfs dev /dev

# Install symlinks 
sudo chroot "$MNT" /bin/busybox --install -s

# Cleanup mount & loop device
sudo umount "$MNT"
sudo losetup -d "$DEVICE"

# 14️⃣ Run QEMU with kernel and rootfs
# Use console for output
qemu-system-x86_64 \
    -kernel "$KERNEL" \
    -append "root=/dev/sda1 init=/init console=ttyS0" \
    -drive format=raw,file="$IMG" \
    -nographic -m 512M