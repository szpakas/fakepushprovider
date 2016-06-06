# generator

Command generates JSON encoded files with applications and instances (devices).
Files are intended to be shared between fake push server and client service under test.

## Usage

    Usage: generator [flags] <apns|fcm> path/to/apps.json path/to/instances.json
    
    flags:
	  -a int
	      number of apps to generate (default: 4)
	  -i int 
	      number of instances per app to generate (default: 150)
	  -u float
	      percent of instances with unregistered status (default: 10.0)
	  -r int
	      maximum number of registrationIDs per app (android only) (default: 10)
    
    examples:
      ./generator -a 8 -i 17 -u 25.4 -r 6 fcm data/apps.json data/ins.json
      ./generator -a 2 -i 4 -u 25 apns tmp/apps.json tmp/ins.json

## Export sample

### FCM/GCM

apps.json

    {"ID":"appId-1","SenderID":"senderId-1","ApiKey":"apiKey-1"}
    {"ID":"appId-2","SenderID":"senderId-2","ApiKey":"apiKey-2"}

ins.json

    {"ID":"appId-1-instanceId-1","State":1,"RegistrationIDS":["appId-1-instanceId-1-regId-1"],"CanonicalID":"appId-1-instanceId-1-regId-1","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-2","State":1,"RegistrationIDS":["appId-1-instanceId-2-regId-1"],"CanonicalID":"appId-1-instanceId-2-regId-1","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-3","State":1,"RegistrationIDS":["appId-1-instanceId-3-regId-1"],"CanonicalID":"appId-1-instanceId-3-regId-1","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-4","State":2,"RegistrationIDS":["appId-1-instanceId-4-regId-1"],"CanonicalID":"appId-1-instanceId-4-regId-1","AppID":"appId-1"}
    {"ID":"appId-2-instanceId-1","State":1,"RegistrationIDS":["appId-2-instanceId-1-regId-1"],"CanonicalID":"appId-2-instanceId-1-regId-1","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-2","State":1,"RegistrationIDS":["appId-2-instanceId-2-regId-1"],"CanonicalID":"appId-2-instanceId-2-regId-1","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-3","State":1,"RegistrationIDS":["appId-2-instanceId-3-regId-1"],"CanonicalID":"appId-2-instanceId-3-regId-1","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-4","State":2,"RegistrationIDS":["appId-2-instanceId-4-regId-1"],"CanonicalID":"appId-2-instanceId-4-regId-1","AppID":"appId-2"}

### APNS

apps.json

    {"ID":"appId-1","BundleID":"bundleId-1"}
    {"ID":"appId-2","BundleID":"bundleId-2"}

ins.json

    {"ID":"appId-1-instanceId-1","State":1,"Token":"appId-1-instanceId-1-token","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-2","State":1,"Token":"appId-1-instanceId-2-token","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-3","State":1,"Token":"appId-1-instanceId-3-token","AppID":"appId-1"}
    {"ID":"appId-1-instanceId-4","State":2,"Token":"appId-1-instanceId-4-token","LastSeen":1465125660,"AppID":"appId-1"}
    {"ID":"appId-2-instanceId-1","State":1,"Token":"appId-2-instanceId-1-token","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-2","State":1,"Token":"appId-2-instanceId-2-token","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-3","State":1,"Token":"appId-2-instanceId-3-token","AppID":"appId-2"}
    {"ID":"appId-2-instanceId-4","State":2,"Token":"appId-2-instanceId-4-token","LastSeen":1465150860,"AppID":"appId-2"}
