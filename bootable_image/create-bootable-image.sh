#!/bin/bash
set -e # Exit if any part fails 

# Temporary workspace 
WORKDIR="$(pwd)/workdir"
mkdir -p "$WORKDIR"

# Copy host kernel to use in QEMU
KERNEL="/boot/vmlinuz-$(uname -r)"
sudo cp "$KERNEL" "$WORKDIR/vmlinuz"
sudo chmod 644 "$WORKDIR/vmlinuz"

# Download static BusyBox so essential Linux commands are accessible 
BUSYBOX="$WORKDIR/busybox"
if [ ! -f "$BUSYBOX" ]; then
    wget -O "$BUSYBOX" https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox
    chmod +x "$BUSYBOX"
fi

# Create initramfs structure
INITRAMFS="$WORKDIR/initramfs"
mkdir -p "$INITRAMFS"/{bin,sbin,etc,proc,sys,usr/bin,usr/sbin,dev}

# Create basic device nodes so messages can be printed 
sudo mknod -m 622 "$INITRAMFS/dev/console" c 5 1
sudo mknod -m 666 "$INITRAMFS/dev/null" c 1 3

# Copy BusyBox as busybox and symlink commands
sudo cp "$BUSYBOX" "$INITRAMFS/bin/busybox"
pushd "$INITRAMFS/bin" > /dev/null
for cmd in sh mount echo cat ls; do
    ln -sf busybox "$cmd"
done
popd > /dev/null

# Create /init script
cat > "$INITRAMFS/init" <<'EOF'
#!/bin/sh
mount -t proc none /proc
mount -t sysfs none /sys
echo "Hello World!" >/dev/console
exec /bin/sh >/dev/console 2>/dev/null
EOF
chmod +x "$INITRAMFS/init"

# Build initramfs 
pushd "$INITRAMFS" > /dev/null
find . | cpio -H newc -o | gzip > "$WORKDIR/initrd.img"
popd > /dev/null

# Create empty disk image 
IMG="$WORKDIR/hello_canonical.img"
qemu-img create -f raw "$IMG" 512M

# Format it and copy the initramfs content
LOOP=$(sudo losetup -f --show "$IMG")
sudo mkfs.ext4 "$LOOP"
MNT="$WORKDIR/mnt"
mkdir -p "$MNT"
sudo mount "$LOOP" "$MNT"
sudo cp -r "$INITRAMFS"/* "$MNT"/
sudo umount "$MNT"
sudo losetup -d "$LOOP"

# Boot QEMU from disk image
qemu-system-x86_64 \
    -m 512M \
    -kernel "$WORKDIR/vmlinuz" \
    -append "root=/dev/sda console=tty0 quiet init=/init rw" \
    -drive file="$IMG",format=raw,if=ide \
    -serial stdio

mv "$IMG" .
rm -rf "$WORKDIR" 

