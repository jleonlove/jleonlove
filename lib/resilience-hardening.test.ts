import {describe,it,expect} from 'vitest';
import {validateWorkflowEvents,validateConcurrentSteps,requireCompensationAfterPartialFailure} from './resilience-hardening';
import type {WorkflowStep} from './workflow-recovery';

describe('RC-000108 resilience hardening',()=>{
 it('fails closed on event sequence corruption',()=>{
  const r=validateWorkflowEvents([
   {id:'e1',workflowId:'w',sequence:1,idempotencyKey:'k1',kind:'START'},
   {id:'e3',workflowId:'w',sequence:3,idempotencyKey:'k3',kind:'SETTLE'}]);
  expect(r.safe).toBe(false); expect(r.findings.some(x=>x.code==='EVENT_SEQUENCE_GAP')).toBe(true);
 });
 it('detects duplicate idempotency events',()=>{
  const r=validateWorkflowEvents([
   {id:'e1',workflowId:'w',sequence:1,idempotencyKey:'same',kind:'START'},
   {id:'e2',workflowId:'w',sequence:2,idempotencyKey:'same',kind:'PAY'}]);
  expect(r.safe).toBe(false);
 });
 it('detects concurrent side-effect race',()=>{
  const steps:WorkflowStep[]=[
   {id:'a',status:'RUNNING',attempts:1,maxAttempts:3,idempotencyKey:'settle-1',compensatable:true},
   {id:'b',status:'RUNNING',attempts:1,maxAttempts:3,idempotencyKey:'settle-1',compensatable:true}];
  expect(validateConcurrentSteps(steps).safe).toBe(false);
 });
 it('requires compensation after downstream partial failure',()=>{
  const steps:WorkflowStep[]=[
   {id:'reserve',status:'SUCCEEDED',attempts:1,maxAttempts:3,idempotencyKey:'r',compensatable:true},
   {id:'transfer',status:'FAILED',attempts:3,maxAttempts:3,idempotencyKey:'t',compensatable:false}];
  expect(requireCompensationAfterPartialFailure(steps).map(x=>x.detail)).toContain('reserve');
 });
});
