package pihole

type DNS struct {
	Upstreams           []string       `json:"upstreams"`
	CNAMEdeepInspect    bool           `json:"CNAMEdeepInspect"`
	BlockESNI           bool           `json:"blockESNI"`
	EDNS0ECS            bool           `json:"EDNS0ECS"`
	IgnoreLocalhost     bool           `json:"ignoreLocalhost"`
	ShowDNSSEC          bool           `json:"showDNSSEC"`
	AnalyzeOnlyAandAAAA bool           `json:"analyzeOnlyAandAAAA"`
	PiholePTR           string         `json:"piholePTR"`
	ReplyWhenBusy       string         `json:"replyWhenBusy"`
	BlockTTL            int            `json:"blockTTL"`
	Hosts               []string       `json:"hosts"`
	DomainNeeded        bool           `json:"domainNeeded"`
	ExpandHosts         bool           `json:"expandHosts"`
	Domain              string         `json:"domain"`
	BogusPriv           bool           `json:"bogusPriv"`
	DNSSEC              bool           `json:"dnssec"`
	Interface           string         `json:"interface"`
	HostRecord          string         `json:"hostRecord"`
	ListeningMode       string         `json:"listeningMode"`
	QueryLogging        bool           `json:"queryLogging"`
	CnameRecords        []string       `json:"cnameRecords"`
	Port                int            `json:"port"`
	RevServers          []string       `json:"revServers"`
	Cache               Cache          `json:"cache"`
	Blocking            Blocking       `json:"blocking"`
	SpecialDomains      SpecialDomains `json:"specialDomains"`
	Reply               Reply          `json:"reply"`
	RateLimit           RateLimit      `json:"rateLimit"`
}

type Cache struct {
	Size               int `json:"size"`
	Optimizer          int `json:"optimizer"`
	UpstreamBlockedTTL int `json:"upstreamBlockedTTL"`
}

type Blocking struct {
	Active bool   `json:"active"`
	Mode   string `json:"mode"`
	Edns   string `json:"edns"`
}

type SpecialDomains struct {
	MozillaCanary      bool `json:"mozillaCanary"`
	ICloudPrivateRelay bool `json:"iCloudPrivateRelay"`
}

type Reply struct {
	Host     Host     `json:"host"`
	Blocking Blocking `json:"blocking"`
}

type Host struct {
	Force4 bool   `json:"force4"`
	IPv4   string `json:"IPv4"`
	Force6 bool   `json:"force6"`
	IPv6   string `json:"IPv6"`
}

type RateLimit struct {
	Count    int `json:"count"`
	Interval int `json:"interval"`
}

type DHCP struct {
	Active               bool     `json:"active"`
	Start                string   `json:"start"`
	End                  string   `json:"end"`
	Router               string   `json:"router"`
	Netmask              string   `json:"netmask"`
	LeaseTime            string   `json:"leaseTime"`
	IPv6                 bool     `json:"ipv6"`
	RapidCommit          bool     `json:"rapidCommit"`
	MultiDNS             bool     `json:"multiDNS"`
	Logging              bool     `json:"logging"`
	IgnoreUnknownClients bool     `json:"ignoreUnknownClients"`
	Hosts                []string `json:"hosts"`
}

type NTP struct {
	IPv4 NTPDetails `json:"ipv4"`
	IPv6 NTPDetails `json:"ipv6"`
	Sync Sync       `json:"sync"`
}

type NTPDetails struct {
	Active  bool   `json:"active"`
	Address string `json:"address"`
}

type Sync struct {
	Active   bool   `json:"active"`
	Server   string `json:"server"`
	Interval int    `json:"interval"`
	Count    int    `json:"count"`
	RTC      RTC    `json:"rtc"`
}

type RTC struct {
	Set    bool   `json:"set"`
	Device string `json:"device"`
	UTC    bool   `json:"utc"`
}

type Resolver struct {
	ResolveIPv4  bool   `json:"resolveIPv4"`
	ResolveIPv6  bool   `json:"resolveIPv6"`
	NetworkNames bool   `json:"networkNames"`
	RefreshNames string `json:"refreshNames"`
}

type Database struct {
	DBImport   bool    `json:"DBimport"`
	MaxDBDays  int     `json:"maxDBdays"`
	DBInterval int     `json:"DBinterval"`
	UseWAL     bool    `json:"useWAL"`
	Network    Network `json:"network"`
}

type Network struct {
	ParseARPcache bool `json:"parseARPcache"`
	Expire        int  `json:"expire"`
}

type Webserver struct {
	Domain    string    `json:"domain"`
	ACL       string    `json:"acl"`
	Port      string    `json:"port"`
	Threads   int       `json:"threads"`
	Headers   []string  `json:"headers"`
	Session   Session   `json:"session"`
	TLS       TLS       `json:"tls"`
	Paths     Paths     `json:"paths"`
	Interface Interface `json:"interface"`
	API       API       `json:"api"`
}

type Session struct {
	Timeout int  `json:"timeout"`
	Restore bool `json:"restore"`
}

type TLS struct {
	Cert string `json:"cert"`
}

type Paths struct {
	Webroot string `json:"webroot"`
	Webhome string `json:"webhome"`
}

type Interface struct {
	Boxed bool   `json:"boxed"`
	Theme string `json:"theme"`
}

type API struct {
	MaxSessions            int      `json:"max_sessions"`
	PrettyJSON             bool     `json:"prettyJSON"`
	PwHash                 string   `json:"pwhash"`
	Password               string   `json:"password"`
	TOTPSecret             string   `json:"totp_secret"`
	AppPwHash              string   `json:"app_pwhash"`
	AppSudo                bool     `json:"app_sudo"`
	CLIPw                  bool     `json:"cli_pw"`
	ExcludeClients         []string `json:"excludeClients"`
	ExcludeDomains         []string `json:"excludeDomains"`
	MaxHistory             int      `json:"maxHistory"`
	MaxClients             int      `json:"maxClients"`
	ClientHistoryGlobalMax bool     `json:"client_history_global_max"`
	AllowDestructive       bool     `json:"allow_destructive"`
	Temp                   Temp     `json:"temp"`
}

type Temp struct {
	Limit int    `json:"limit"`
	Unit  string `json:"unit"`
}

type Files struct {
	PID        string `json:"pid"`
	Database   string `json:"database"`
	Gravity    string `json:"gravity"`
	GravityTmp string `json:"gravity_tmp"`
	MacVendor  string `json:"macvendor"`
	SetupVars  string `json:"setupVars"`
	PCAP       string `json:"pcap"`
	Log        Log    `json:"log"`
}

type Log struct {
	FTL       string `json:"ftl"`
	Dnsmasq   string `json:"dnsmasq"`
	Webserver string `json:"webserver"`
}

type Misc struct {
	PrivacyLevel int      `json:"privacylevel"`
	DelayStartup int      `json:"delay_startup"`
	Nice         int      `json:"nice"`
	Addr2line    bool     `json:"addr2line"`
	EtcDnsmasqD  bool     `json:"etc_dnsmasq_d"`
	DnsmasqLines []string `json:"dnsmasq_lines"`
	ExtraLogging bool     `json:"extraLogging"`
	ReadOnly     bool     `json:"readOnly"`
	Check        Check    `json:"check"`
}

type Check struct {
	Load  bool `json:"load"`
	Shmem int  `json:"shmem"`
	Disk  int  `json:"disk"`
}

type Debug struct {
	Database     bool `json:"database"`
	Networking   bool `json:"networking"`
	Locks        bool `json:"locks"`
	Queries      bool `json:"queries"`
	Flags        bool `json:"flags"`
	Shmem        bool `json:"shmem"`
	GC           bool `json:"gc"`
	ARP          bool `json:"arp"`
	Regex        bool `json:"regex"`
	API          bool `json:"api"`
	TLS          bool `json:"tls"`
	Overtime     bool `json:"overtime"`
	Status       bool `json:"status"`
	Caps         bool `json:"caps"`
	DNSSEC       bool `json:"dnssec"`
	Vectors      bool `json:"vectors"`
	Resolver     bool `json:"resolver"`
	EDNS0        bool `json:"edns0"`
	Clients      bool `json:"clients"`
	AliasClients bool `json:"aliasclients"`
	Events       bool `json:"events"`
	Helper       bool `json:"helper"`
	Config       bool `json:"config"`
	Inotify      bool `json:"inotify"`
	Webserver    bool `json:"webserver"`
	Extra        bool `json:"extra"`
	Reserved     bool `json:"reserved"`
	NTP          bool `json:"ntp"`
	Netlink      bool `json:"netlink"`
	All          bool `json:"all"`
}

type Config struct {
	DNS       DNS       `json:"dns"`
	DHCP      DHCP      `json:"dhcp"`
	NTP       NTP       `json:"ntp"`
	Resolver  Resolver  `json:"resolver"`
	Database  Database  `json:"database"`
	Webserver Webserver `json:"webserver"`
	Files     Files     `json:"files"`
	Misc      Misc      `json:"misc"`
	Debug     Debug     `json:"debug"`
}

type ConfigResponse struct {
	Config Config  `json:"config"`
	Took   float64 `json:"took"`
}
