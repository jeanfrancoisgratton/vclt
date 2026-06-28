%define debug_package   %{nil}
%define _build_id_links none
%define _name vclt
%define _prefix /opt
%define _bash_completionsdir /usr/share/bash-completion/completions
%define _zsh_completionsdir  /usr/share/zsh/site-functions
%define _version 2.3.2
%define _rel 0
%define _binaryname vclt

Name:       vclt
Version:    %{_version}
Release:    %{_rel}
Summary:    Hashicorp Vault client

Group:      CI/CD
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/devops/vclt.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc
Recommends: zsh
Requires: bash-completion

%description
Hashicorp Vault client

%prep
%autosetup

%build
cd src
go mod download
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_builddir}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre
%install
install -Dpm 0755 %{_builddir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
# Bash completion — always install
/opt/bin/vclt completion bash > %{_bash_completionsdir}/vclt

# Zsh completion — only if zsh is present
if command -v zsh > /dev/null 2>&1; then
    mkdir -p %{_zsh_completionsdir}/zsh/site-functions
    /opt/bin/vclt completion zsh > %{_zsh_completionsdir}/_vclt
fi

%preun

%postun
if [ $1 -eq 0 ]; then
    # $1 == 0 means this is a full uninstall, not an upgrade
    rm -f %{_bash_completionsdir}/vclt
    rm -f %{_zsh_completionsdir}/_vclt
fi

%files
%defattr(0755,root,root,-)
%{_bindir}/%{_binaryname}

%changelog
* Sun Jun 28 2026 Binary package builder <builder@famillegratton.net> 2.3.2-0
- fixed issue where the temp files were left on the filesystem
- chore: update changelog for 2.3.1-0
- Backup and restored kv engines now encode their JSON files
- fixed path to created binary in DEB building package
- Fixed debian build script

* Sun Jun 28 2026 Binary package builder <builder@famillegratton.net> 2.3.1-0
- Backup and restored kv engines now encode their JSON files
- fixed path to created binary in DEB building package
- Fixed debian build script

* Sat Jun 27 2026 Binary package builder <builder@famillegratton.net> 2.3.0-0
- fixed debian build script
- completed backup and restore
- fixed go vet issue with Printf/Print
- chore: update changelog for 2.2.0-0
- version bump
- Potential fix for code scanning alert no. 1: Incorrect conversion between integer types
- fixed TUI when prompting for a disable confirmation
- Completed KV engine enable/disable
- moved secrets to kv, created new sys package
- commit before major refactoring
- admin setrootkeys can now work in offline mode
- Added the initial root key field to rootkeys json
- completed policy write and read
- Version numbering scheme change to align with SemVer
- added the list mounts command
- added the list mounts command
- completed policies rm
- completed policies read
- initial stub for policies
- fixed binary name in DEBIAN build scripts

* Fri Jun 26 2026 Binary package builder <builder@famillegratton.net> 2.2.0-0
- version bump
- Potential fix for code scanning alert no. 1: Incorrect conversion between integer types
- fixed TUI when prompting for a disable confirmation
- Completed KV engine enable/disable
- moved secrets to kv, created new sys package
- commit before major refactoring
- admin setrootkeys can now work in offline mode
- Added the initial root key field to rootkeys json
- completed policy write and read
- Version numbering scheme change to align with SemVer
- added the list mounts command
- added the list mounts command
- completed policies rm
- completed policies read
- initial stub for policies
- fixed binary name in DEBIAN build scripts

* Thu Jun 18 2026 Binary package builder <builder@famillegratton.net> 2.00.00~DEBUG3-0
- Completed secrets rm and secrets destroy
- chore: update changelog for 2.00.00~DEBUG2-0
- completed secrets put
- fixed archlinux erroneous revision number
- chore: update changelog for 2.00.00~DEBUG-0
- completed secrets write
- rewrote specfile and Makefile
- build script fixes for rhel
- changed paths in buildeps script
- Fixed makefile
- fixed all packaging scripts
- fixed numeral outputting
- -x needs fixing
- removed left over files
- removed ENV; secrets read needs testing
- completed read subcommand
- more cleanup before starting afresh
- Stub commit
- interim sync across workspaces

* Wed Jun 17 2026 Binary package builder <builder@famillegratton.net> 2.00.00~DEBUG2-0
- completed secrets put
- fixed archlinux erroneous revision number
- chore: update changelog for 2.00.00~DEBUG-0
- completed secrets write
- rewrote specfile and Makefile
- build script fixes for rhel
- changed paths in buildeps script
- Fixed makefile
- fixed all packaging scripts
- fixed numeral outputting
- -x needs fixing
- removed left over files
- removed ENV; secrets read needs testing
- completed read subcommand
- more cleanup before starting afresh
- Stub commit
- interim sync across workspaces

* Tue Jun 16 2026 Binary package builder <builder@famillegratton.net> 2.00.00~DEBUG-0
- completed secrets write
- rewrote specfile and Makefile
- build script fixes for rhel
- changed paths in buildeps script
- Fixed makefile
- fixed all packaging scripts
- fixed numeral outputting
- -x needs fixing
- removed left over files
- removed ENV; secrets read needs testing
- completed read subcommand
- more cleanup before starting afresh
- Stub commit
- interim sync across workspaces

* Mon Jun 15 2026 Binary package builder <builder@famillegratton.net> 2.00.00-0
- build script fixes for rhel
- changed paths in buildeps script
- Fixed makefile
- fixed all packaging scripts
- fixed numeral outputting
- -x needs fixing
- removed left over files
- removed ENV; secrets read needs testing
- completed read subcommand
- more cleanup before starting afresh
- Stub commit
- interim sync across workspaces

* Wed Jul 03 2024 RPM Builder <builder@famillegratton.net> 1.01.00-1
- Version bump as -v was not being correctly displayed (jean-
  francois@famillegratton.net)

* Tue Jul 02 2024 RPM Builder <builder@famillegratton.net> 1.01.00-0
- Version bump pursuant previous commit (jean-francois@famillegratton.net)
- Fixed some security vulns in dependencies (jean-francois@famillegratton.net)

* Tue Jul 02 2024 RPM Builder <builder@famillegratton.net> 1.00.00-0
- new package built with tito

