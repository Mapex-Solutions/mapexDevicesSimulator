#!/bin/bash
# This is the COMPLETE .deb postinst — electron-builder's `deb.afterInstall`
# replaces the generated script rather than appending, so it reproduces the
# default integration (CLI symlink, mime/desktop DB refresh) and, crucially,
# forces the SUID chrome-sandbox unconditionally.
#
# Install dir has NO spaces (/opt/MapexDeviceSimulator) — spaces break Chromium's
# zygote/sandbox child launch on Linux.

if type update-alternatives 2>/dev/null >&1; then
    # Remove previous link if it doesn't use update-alternatives
    if [ -L '/usr/bin/mapex-devices-simulator' -a -e '/usr/bin/mapex-devices-simulator' -a "`readlink '/usr/bin/mapex-devices-simulator'`" != '/etc/alternatives/mapex-devices-simulator' ]; then
        rm -f '/usr/bin/mapex-devices-simulator'
    fi
    update-alternatives --install '/usr/bin/mapex-devices-simulator' 'mapex-devices-simulator' '/opt/MapexDeviceSimulator/mapex-devices-simulator' 100 || ln -sf '/opt/MapexDeviceSimulator/mapex-devices-simulator' '/usr/bin/mapex-devices-simulator'
else
    ln -sf '/opt/MapexDeviceSimulator/mapex-devices-simulator' '/usr/bin/mapex-devices-simulator'
fi

# Force the SUID chrome-sandbox so the app runs WITH the sandbox on systems where
# the unprivileged user-namespace sandbox is restricted — notably Ubuntu 24.04,
# whose AppArmor blocks the userns operations Electron's zygote needs, making it
# fall back to the setuid sandbox. 4755 makes that path work, no --no-sandbox.
chmod 4755 '/opt/MapexDeviceSimulator/chrome-sandbox' || true

if hash update-mime-database 2>/dev/null; then
    update-mime-database /usr/share/mime || true
fi

if hash update-desktop-database 2>/dev/null; then
    update-desktop-database /usr/share/applications || true
fi
