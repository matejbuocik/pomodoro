# https://github.com/matejbuocik/pomodoro
%global goipath         github.com/matejbuocik/pomodoro
Version: v0.1

%global common_description %{expand:
TUI Pomodoro Timer 🍅. Written in Go.
}

%gometa

%global golicenses      LICENSE
%global godocs          README.md

Name:           %{goname}
Release:        1%{?dist}
Summary:        TUI Pomodoro Timer 🍅

License:        MIT
URL:            %{gourl}
Source:         %{gosource}

%description
%{common_description}

%gopkg

%prep
%goprep

%generate_buildrequires
%go_generate_buildrequires

%build
%gobuild -o %{gobuilddir}/bin/pomodoro %{goipath}

%install
install -m 0755 -vd                     %{buildroot}%{_bindir}
install -m 0755 -vp %{gobuilddir}/bin/* %{buildroot}%{_bindir}/

%check
%gocheck

%files
%license LICENSE
%doc README.md
%{_bindir}/*

%changelog
* Mon Sep 23 2024 15:00:00 CET Matej Buocik <matej.buocik@gmail.com> - v0.1
- First release
