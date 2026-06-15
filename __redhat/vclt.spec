%define debug_package   %{nil}
%define _build_id_links none
%define _name vclt
%define _prefix /opt
%define _version 2.00.00
%define _rel 0
%define _arch x86_64
%define _binaryname vclt

Name:       vclt
Version:    %{_version}
Release:    %{_rel}
Summary:    Hashicorp Vault client

Group:      CI/CD utils
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/devops/vclt

Source0:    %{name}-%{_version}.tar.gz
BuildArchitectures: x86_64
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
Hashicorp Vault client

%prep
%autosetup

%build
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin go build -o %{_sourcedir}/%{_binaryname} .
strip %{_sourcedir}/%{_binaryname}

%clean
rm -rf $RPM_BUILD_ROOT

%pre
exit 0

%install
install -Dpm 0755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
# Bash completion — always install
/opt/bin/vclt completion bash > %{_bash_completionsdir}/vclt

# Zsh completion — only if zsh is present
if command -v zsh > /dev/null 2>&1; then
    mkdir -p /usr/share/zsh/site-functions
    /opt/bin/vclt completion zsh > /usr/share/zsh/site-functions/_vclt
    zsh -c 'autoload -Uz compinit && compinit' 2>/dev/null || true
fi

%preun

%postun
if [ $1 -eq 0 ]; then
    # $1 == 0 means this is a full uninstall, not an upgrade
    rm -f %{_bash_completionsdir}/vclt
    rm -f %{_zsh_completionsdir}/_vclt
fi

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Wed Jul 03 2024 RPM Builder <builder@famillegratton.net> 1.01.00-1
- Version bump as -v was not being correctly displayed (jean-
  francois@famillegratton.net)

* Tue Jul 02 2024 RPM Builder <builder@famillegratton.net> 1.01.00-0
- Version bump pursuant previous commit (jean-francois@famillegratton.net)
- Fixed some security vulns in dependencies (jean-francois@famillegratton.net)

* Tue Jul 02 2024 RPM Builder <builder@famillegratton.net> 1.00.00-0
- new package built with tito

