Notify
======
[![Build Status](https://travis-ci.org/Danzabar/notify.svg?branch=master)](https://travis-ci.org/Danzabar/notify) [![Coverage Status](https://coveralls.io/repos/github/Danzabar/notify/badge.svg?branch=master)](https://coveralls.io/github/Danzabar/notify?branch=master)

A micro service to store and serve notifications from various sources

## Usage
First install the package

	go install github.com/Danzabar/notify

Then run the application with your settings, by default this will use a sqlite connection:

	notify -driver="mysql" -creds="user:pass@/dbname?charset=utf8"

Use the `-m` flag to run migrations

	notify -m

### Sending Alerts

To start scanning and sending alerts use the `-a` flag

	notify -a
