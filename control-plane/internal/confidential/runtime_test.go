package confidential
import("errors";"testing")
func pol()Policy{return Policy{Regions:map[string]bool{"us":true},OfflineCapabilities:map[string]bool{"read":true,"infer":true}}}
func TestEdgeSensitiveAllowed(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Edge,Region:"us",Sensitive:true,Capability:"infer"});e!=nil{t.Fatal(e)}}
func TestSensitiveCloudDenied(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Cloud,Region:"us",Sensitive:true});!errors.Is(e,ErrPlacement){t.Fatal(e)}}
func TestConfidentialNeedsAttestation(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Confidential,Region:"us",Sensitive:true});!errors.Is(e,ErrAttestation){t.Fatal(e)}}
func TestAttestedConfidentialAllowed(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Confidential,Region:"us",Sensitive:true,Attested:true});e!=nil{t.Fatal(e)}}
func TestResidencyDenied(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Edge,Region:"eu"});!errors.Is(e,ErrResidency){t.Fatal(e)}}
func TestOfflineWriteDenied(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Edge,Region:"us",Offline:true,Capability:"write"});!errors.Is(e,ErrOffline){t.Fatal(e)}}
func TestOfflineInferenceAllowed(t *testing.T){if e:=Authorize(pol(),Workload{Mode:Edge,Region:"us",Offline:true,Capability:"infer"});e!=nil{t.Fatal(e)}}
