package governance

import (
	"context"
	"errors"
)

var ErrDenied = errors.New("governance denied execution")

type Request struct {
	RequestID string
	AgentID string
	TaskID string
	ReleaseID string
	CapabilityID string
	ScopeID string
	TrajectoryID string
	Resource string
	Action string
	Arguments map[string]any
}

type Result struct { Data []byte }

type Authority interface { Verify(context.Context,Request) error }
type Scope interface { Enforce(context.Context,Request) error }
type Policy interface { Evaluate(context.Context,Request) (bool,error) }
type Trajectory interface { Evaluate(context.Context,Request) error }
type Containment interface { Attest(context.Context,Request) error }
type Evidence interface {
	Intent(context.Context,Request) error
	Denial(context.Context,Request,string) error
	Outcome(context.Context,Request,Result,error) error
}
type Runtime interface { Execute(context.Context,Request)(Result,error) }

type Pipeline struct {
	Authority Authority
	Scope Scope
	Policy Policy
	Trajectory Trajectory
	Containment Containment
	Evidence Evidence
	Runtime Runtime
}

func (p *Pipeline) Execute(ctx context.Context,r Request)(Result,error){
	deny:=func(reason string,cause error)(Result,error){
		_ = p.Evidence.Denial(ctx,r,reason)
		return Result{},errors.Join(ErrDenied,cause)
	}
	if err:=p.Authority.Verify(ctx,r);err!=nil{return deny("authority",err)}
	if err:=p.Scope.Enforce(ctx,r);err!=nil{return deny("scope",err)}
	if err:=p.Trajectory.Evaluate(ctx,r);err!=nil{return deny("trajectory",err)}
	if err:=p.Containment.Attest(ctx,r);err!=nil{return deny("containment",err)}
	ok,err:=p.Policy.Evaluate(ctx,r)
	if err!=nil{return deny("policy_unavailable",err)}
	if !ok{return deny("policy_denied",ErrDenied)}
	if err:=p.Evidence.Intent(ctx,r);err!=nil{return Result{},errors.Join(ErrDenied,err)}
	out,execErr:=p.Runtime.Execute(ctx,r)
	if err:=p.Evidence.Outcome(ctx,r,out,execErr);err!=nil{return Result{},errors.Join(execErr,err)}
	return out,execErr
}
