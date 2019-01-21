%define debug_package %{nil}

Name:    ves-agent
Version: %{getenv:VERSION}
# Release: 1%{?dist}
Release: %{getenv:BUILDID}
Summary: Bridge between prometheus and ONAP's VES-Collector.
License: Copyright Nokia 2018

Source0: ves-agent
Source1: ves-agent.service
Source2: ves-agent.yml

%{?systemd_requires}
Requires(pre): shadow-utils

%description

VES-Agent is a service acting as a bridge between prometheus and ONAP's VES-Collector.

%build
/bin/true

%install
mkdir -vp %{buildroot}%{_sharedstatedir}/ves-agent
install -D -m 755 %{SOURCE0} %{buildroot}%{_bindir}/ves-agent
install -D -m 644 %{SOURCE1} %{buildroot}%{_unitdir}/ves-agent.service
install -D -m 644 %{SOURCE2} %{buildroot}%{_sysconfdir}/ves-agent/ves-agent.yml

%pre
getent group prometheus >/dev/null || groupadd -r prometheus
getent passwd prometheus >/dev/null || \
  useradd -r -g prometheus -d %{_sharedstatedir}/prometheus -s /sbin/nologin \
          -c "Prometheus services" prometheus
exit 0

%post
%systemd_post ves-agent.service

%preun
%systemd_preun ves-agent.service

%postun
%systemd_postun ves-agent.service

%files
%defattr(-,root,root,-)
%{_bindir}/ves-agent
%{_unitdir}/ves-agent.service
%config(noreplace) %{_sysconfdir}/ves-agent/ves-agent.yml
%dir %attr(755, prometheus, prometheus)%{_sharedstatedir}/ves-agent
