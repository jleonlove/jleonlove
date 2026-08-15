import {describe,it,expect} from "vitest";import {fastestLegitimatePathToClose} from "./fastest-path-to-close";
describe("fastest legitimate path to close",()=>{
 it("prioritizes a hard blocker that unlocks downstream work",()=>{const p=fastestLegitimatePathToClose([{id:"title",label:"Verify title",owner:"seller",dependsOn:[],state:"OPEN",hardBlocker:true,unlocks:["close"]},{id:"format",label:"Format memo",owner:"ops",dependsOn:[],state:"OPEN",unlocks:[]}]);expect(p.nextAction).toBe("title")});
 it("keeps independent work parallel",()=>{const p=fastestLegitimatePathToClose([{id:"inspection",label:"Inspection",owner:"inspector",dependsOn:[],state:"BLOCKED"},{id:"bank",label:"Bank verification",owner:"buyer",dependsOn:[],state:"OPEN"}]);expect(p.parallelNow).toContain("bank")});
 it("detects circular dependency deadlocks",()=>{const p=fastestLegitimatePathToClose([{id:"a",label:"A",owner:"seller",dependsOn:["b"],state:"OPEN"},{id:"b",label:"B",owner:"bank",dependsOn:["a"],state:"OPEN"}]);expect(p.deadlocks.length).toBeGreaterThan(0)});
});
