#!/usr/bin/env sh

export {{.AppIdName}}="{{.AppIdValue}}"
export {{.AppKeyName}}="{{.AppKeyValue}}"

. "/root/.acme.sh/acme.sh.env"

acme.sh --upgrade

acme.sh --register-account -m admin@lingyin99.cn --issue --dns {{.DnsApi}} -d {{.MainDomain}} {{.ExtraDomain}}