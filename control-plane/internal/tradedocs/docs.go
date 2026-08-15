package tradedocs
import("errors";"sort";"strings")
var(ErrType=errors.New("unsupported trade document");ErrProvenance=errors.New("document provenance required");ErrConflict=errors.New("trade document conflict");ErrField=errors.New("required field missing"))
type Document struct{ID,Type,Source,Hash string;Fields map[string]string}
type Conflict struct{Field string;Values map[string]string}
var required=map[string][]string{"SPA":{"commodity","quantity","unit","buyer","seller"},"INVOICE":{"commodity","quantity","unit"},"BL":{"commodity","quantity","unit","origin","destination"},"COA":{"commodity"}}
func Validate(d Document)error{
 req,ok:=required[d.Type];if !ok{return ErrType}
 if strings.TrimSpace(d.Source)==""||strings.TrimSpace(d.Hash)==""{return ErrProvenance}
 for _,f:=range req{if strings.TrimSpace(d.Fields[f])==""{return ErrField}}
 return nil
}
func Reconcile(ds []Document,fields []string)([]Conflict,error){
 for _,d:=range ds{if e:=Validate(d);e!=nil{return nil,e}}
 var out []Conflict
 for _,f:=range fields{vals:=map[string]string{};uniq:=map[string]bool{};for _,d:=range ds{if v:=strings.TrimSpace(d.Fields[f]);v!=""{vals[d.ID]=v;uniq[strings.ToLower(v)]=true}}
  if len(uniq)>1{out=append(out,Conflict{Field:f,Values:vals})}}
 sort.Slice(out,func(i,j int)bool{return out[i].Field<out[j].Field});return out,nil
}
