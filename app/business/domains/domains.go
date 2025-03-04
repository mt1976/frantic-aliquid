package domains

import "strings"

type Domain struct {
	name   string
	normal bool
}

func (d Domain) String() string {
	return strings.ToTitle(d.name)
}

func (d Domain) Equals(other Domain) bool {
	if d.name == other.name && d.normal == other.normal {
		return true
	}
	return false
}

func (d Domain) IsDefault() bool {
	return d.normal
}

func Special(s string) Domain {
	return Domain{s, false}
}

var TEXT = Domain{"texts", true}
var AUTHORITY = Domain{"authorities", true}
var BEHAVIOUR = Domain{"behaviours", true}
var HOSTS = Domain{"hosts", true}
var PASSWORD = Domain{"passwords", true}
var SESSION = Domain{"sessions", true}
var SETTINGS = Domain{"settings", true}
var STATUS = Domain{"status", true}
var USER = Domain{"users", true}
var ZONE = Domain{"zones", true}
var SECURITY = Domain{"security", true}
var API = Domain{"apis", true}
var ROUTE = Domain{"routes", true}
var DATABASE = Domain{"database", true}
var JOBS = Domain{"jobs", true}
var NOTIFICATIONS = Domain{"notifications", true}
