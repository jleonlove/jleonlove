package qualification
import("strings")
type Policy struct{AllowedHosts map[string]bool;AllowedPaths map[string]bool}
type Probe struct{NetworkHosts,FilesystemPaths,RequestedSecrets []string;Persistence,PrivilegeEscalation bool;InstructionBehavior string}
type Result struct{Qualified bool;Violations []string}
func Evaluate(p Policy,x Probe)Result{var v []string;for _,h:=range x.NetworkHosts{if !p.AllowedHosts[h]{v=append(v,"network")}};for _,f:=range x.FilesystemPaths{if !p.AllowedPaths[f]{v=append(v,"filesystem")}};if len(x.RequestedSecrets)>0{v=append(v,"secret")};if x.Persistence{v=append(v,"persistence")};if x.PrivilegeEscalation{v=append(v,"privilege")};b:=strings.ToLower(x.InstructionBehavior);if strings.Contains(b,"ignore policy")||strings.Contains(b,"override system"){v=append(v,"prompt-injection")};return Result{Qualified:len(v)==0,Violations:v}}
