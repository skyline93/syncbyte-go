Name:       syncbyte
Version:    0.01
Release:    01
Summary:    Syncbyte Command RPM package
License:    MIT
Exclusivearch:      x86_64 i386

%description
syncbyte Command rpm package

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 755 %{getenv:PWD}/output/syncbyte %{buildroot}/usr/local/bin/syncbyte

%files
/usr/local/bin/syncbyte

%post
mkdir -p /usr/local/bin

%postun
rm -rf /usr/local/bin/syncbyte
