Name:       syncbyte-agent
Version:    0.01
Release:    01
Summary:    Syncbyte Agent RPM package
License:    MIT
Exclusivearch:      x86_64 i386

%description
syncbyte agent rpm package

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/syncbyte
mkdir -p %{buildroot}/var/log/syncbyte
mkdir -p %{buildroot}/var/run/backup
mkdir -p %{buildroot}/var/run/restore
mkdir -p %{buildroot}/usr/lib/systemd/system

install -m 660 %{getenv:PWD}/conf/agent.yaml %{buildroot}/etc/syncbyte/agent.yaml
install -m 755 %{getenv:PWD}/output/syncbyte-agent %{buildroot}/usr/local/bin/syncbyte-agent
install -m 644 %{getenv:PWD}/deploy/syncbyte-agent.service %{buildroot}/usr/lib/systemd/system/syncbyte-agent.service

%files
/usr/local/bin/syncbyte-agent
/etc/syncbyte/agent.yaml
/usr/lib/systemd/system/syncbyte-agent.service

%post
mkdir -p /usr/local/bin
mkdir -p /etc/syncbyte
mkdir -p /var/log/syncbyte
mkdir -p /var/run/backup
mkdir -p /var/run/restore
mkdir -p /usr/lib/systemd/system

%postun
rm -rf /usr/local/bin/syncbyte-agent
rm -rf /etc/syncbyte/agent.yaml
rm -rf /usr/lib/systemd/system/syncbyte-agent.service
