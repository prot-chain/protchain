"use strict";

const { WorkloadModuleBase } = require("@hyperledger/caliper-core");
const fs = require('fs');
const path = require('path');

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

        
        // Write the generated data to a JSON file
        const outputPath = path.resolve(__dirname, 'protein_data.json');
        // Ensure the file exists before reading
        if (!fs.existsSync(outputPath)) {
            console.log(`File "${outputPath}" does not exist. Creating a new file...`);
            // Add the new transaction data
            const data = {
                proteinId: proteinId,
                fileURL: randomFileURL,
                hash: randomHash
            }

            // Save the updated data back to the file
            fs.writeFileSync(outputPath, JSON.stringify(data, null, 2));
        }

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

