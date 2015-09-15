# Name Tag Auto Printing

##Printer
<pre>
type Printer struct {
	Id                           uuid.UUID
	Name, Ip, ApiKey, ConfigFile string
	Port                         int
	Active, Printing             bool
	NameTag                      *NameTag
}
</pre>