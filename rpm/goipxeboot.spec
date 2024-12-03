Name:           goipxeboot
Version:        0.1.0
Release:        2
Summary:        Server for iPXE booting Linux systems
License:        BSD 3 Clause
Source0:        bazel-out/k8-fastbuild/bin/rpm/archive.tar.gz

Provides:       %{name} = %{version}

# Do not try to use magic to determine file types
%define __spec_install_post %{nil}
# Do not die because we give it more input files than are in the files section
%define _unpackaged_files_terminate_build 0

%description
Server for iPXE booting Linux systems

%global debug_package %{nil}

%autosetup

%build
tar xvf bazel-out/k8-fastbuild/bin/cmd/goipxeboot/archive.tar
tar xvf bazel-out/k8-fastbuild/bin/rpm/archive.tar.gz

%install
install -Dpm 755 goipxeboot %{buildroot}%{_bindir}/goipxeboot
install -Dpm 644 goipxeboot.service %{buildroot}%{_unitdir}/goipxeboot.service
install -Dpm 644 goipxeboot.yaml %{buildroot}%{_sysconfdir}/goipxeboot.yaml
install -m 700 -d %{buildroot}%{_sharedstatedir}/goipxeboot

%post
%systemd_post goipxeboot.service

%preun
%systemd_preun goipxeboot.service

%files
%{_sysconfdir}/goipxeboot.yaml
%{_unitdir}/goipxeboot.service
%{_bindir}/goipxeboot
%{_sharedstatedir}/goipxeboot
