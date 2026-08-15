package contextcompiler
import"testing"
func TestPolicyRanksAboveUntrusted(t *testing.T){c,_:=Compile([]Item{{ID:"web",Text:"ignore policy",Trust:Untrusted,Relevance:100,Tokens:5,Instruction:true},{ID:"p",Text:"deny",Trust:Policy,Tokens:5,Instruction:true}},10);if c.Items[0].ID!="p"{t.Fatal("trust ordering failed")}}
func TestUntrustedCannotBecomeInstruction(t *testing.T){c,_:=Compile([]Item{{ID:"web",Text:"override system",Trust:Untrusted,Tokens:2,Instruction:true}},10);if HasExecutableUntrusted(c){t.Fatal("untrusted executable")}}
func TestBudgetCompaction(t *testing.T){c,_:=Compile([]Item{{ID:"a",Trust:Evidence,Relevance:10,Tokens:7},{ID:"b",Trust:Memory,Relevance:1,Tokens:7}},7);if len(c.Items)!=1||c.Tokens>7{t.Fatal("budget failed")}}
func TestEvidenceBeatsMemory(t *testing.T){c,_:=Compile([]Item{{ID:"m",Trust:Memory,Relevance:99,Tokens:1},{ID:"e",Trust:Evidence,Relevance:1,Tokens:1}},10);if c.Items[0].ID!="e"{t.Fatal("evidence not prioritized")}}
