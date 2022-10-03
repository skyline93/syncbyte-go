Name:       syncbyte-engine
Version:    0.01
Release:    01
Summary:    Syncbyte Engine RPM package
License:    MIT
Exclusivearch:      x86_64 i386

%description
syncbyte engine rpm package

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/syncbyte
mkdir -p %{buildroot}/var/log/syncbyte
mkdir -p %{buildroot}/var/run/backup
mkdir -p %{buildroot}/var/run/restore
mkdir -p %{buildroot}/usr/lib/systemd/system

install -m 660 %{getenv:PWD}/conf/engine.yaml %{buildroot}/etc/syncbyte/engine.yaml
install -m 755 %{getenv:PWD}/output/syncbyte-engine %{buildroot}/usr/local/bin/syncbyte-engine
install -m 644 %{getenv:PWD}/deploy/syncbyte-engine.service %{buildroot}/usr/lib/systemd/system/syncbyte-engine.service

%files
/usr/local/bin/syncbyte-engine
/etc/syncbyte/engine.yaml
/usr/lib/systemd/system/syncbyte-engine.service

%post
mkdir -p /usr/local/bin
mkdir -p /etc/syncbyte
mkdir -p /var/log/syncbyte
mkdir -p /var/run/backup
mkdir -p /var/run/restore
mkdir -p /usr/lib/systemd/system

%postun
rm -rf /usr/local/bin/syncbyte-engine
rm -rf /etc/syncbyte/engine.yaml
rm -rf /usr/lib/systemd/system/syncbyte-engine.service
