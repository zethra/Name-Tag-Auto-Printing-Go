# Name Tag Auto Printing

##Printer
<pre>
Printer {
	Id			uuid.UUID 	Unique ID of the printer
	Name		string		Readable name of the printer
	Ip			string		The printers ip address
	Port		int			The port octoprint is running on
	ApiKey		string		The printers octoprint api key
	ConfigFile	string		The printers slic3r config file (default used if non provided)
	Active		bool		Whether the printer is useable
	Printing	bool		Whether the printer is printing
	NameTame	*NameTag	A pointer to the name tag currently assign to the printer
}
</pre>

##Name Tag
<pre>
NameTag {
	Id			uuid.UUID 	Unique ID of the name tag
	Name		string		Readable name of the name tag
	Stl			string		The path to the name tag's stl file (blank if not created)
	Gcode		string		The path to the name tag's gcode file (blank if not created)
	Printing	bool		Whether the name tga is printing
	Error		bool		Did the name tag enconter an error (system will ignore name tag if true)
}
</pre>

##API
###Add Name Tag
GET /queue/add
####Form paramitors:
 - **name** - name tag name (reqired) - string
 - **stl** - stl location - string
 - **gcode** - gcode location - string
 - **printing** - printing or not - boolean
 - **error** - error or not - boolean

###Add Printer
GET /printers/add
####Form paramitors:
 - **name** - name tag name (reqired) - string
 - **ip** - printer ip - string
 - **port** - printer port - integer
 - **apiKey** - octoprint api key - string
 - **configFile** - slic3r config file path - string
 - **active** - active or not - boolean
 - **printing** - printing or not - boolean
