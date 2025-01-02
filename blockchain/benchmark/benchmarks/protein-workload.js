"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");

class ProteinWorkload extends WorkloadModuleBase {

    // Called once for each Caliper worker (client process) before the test round
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.workerIndex = workerIndex;
        this.sutAdapter = sutAdapter;
        this.sutContext = sutContext;
        this.txIndex = 0;
    }

    // The function Caliper calls repeatedly to submit transactions
    async submitTransaction() {
        // We'll alternate between storing and querying metadata 
        // (You can separate them into different test rounds if you prefer)

        this.txIndex++;
        const proteinId = `PROT${this.workerIndex}_${this.txIndex}`;

        const randomFileURL = `https://example.com/protein_${this.txIndex}`;
        const randomHash = `hash_${this.txIndex}`;

        // Prepare a random transaction request
        const request = {
            contractId: "proteinmetadata",    // The chaincode name from your Fabric deployment
            contractVersion: "1",        // If you have a version, specify it; else omit
            contractFunction: "StoreMetadata",
            invokerIdentity: "User1",
            contractArguments: [randomHash, proteinId, randomFileURL],
            readOnly: false
        };

        // Submit the transaction
        await this.sutAdapter.sendRequests(request);

        // (Optional) Right after storing, we could query the same item
        // to measure read latency in the same round:
        // const queryRequest = {
        //     contractId: "proteincc",
        //     transaction: "QueryMetadata",
        //     args: [proteinId],
        //     readOnly: true
        // };
        // await this.sutAdapter.sendRequests(queryRequest);
    }
}

function createWorkloadModule() {
    return new ProteinWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;

