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

        // Define the path to the file
        const inputPath = path.resolve(__dirname, 'protein_data.json');

        // Check if the file exists
        if (!fs.existsSync(inputPath)) {
            console.error('The file "protein_data.json" does not exist. Please run the workload to generate it first.');
            return;
        }

        // Read and parse the file content
        const fileContent = fs.readFileSync(inputPath, 'utf8');
        const proteinData = JSON.parse(fileContent);


        // Get the stored data (assuming only the first entry is written)
        const firstEntry = proteinData;
        //console.log(`Retrieved data from file:`);
        //console.log(`Protein ID: ${firstEntry.proteinId}`);
        //console.log(`File URL: ${firstEntry.fileURL}`);
        //console.log(`Hash: ${firstEntry.hash}`);
        const proteinId = firstEntry.proteinId
        //const randomFileURL = `https://example.com/protein_${this.txIndex}`;
        //const randomHash = `hash_${this.txIndex}`;

        // Prepare a random transaction request
        const request = {
            contractId: "proteinmetadata",    // The chaincode name from your Fabric deployment
            contractVersion: "1",        // If you have a version, specify it; else omit
            contractFunction: "QueryMetadata",
            invokerIdentity: "User1",
            contractArguments: [ proteinId ],
            readOnly: false
        };

        let x = 1;
        if (x == 1) {
          console.log("=========== "+proteinId+"==============")
        }

        // Submit the transaction
        await this.sutAdapter.sendRequests(request);


    }
}

function createWorkloadModule() {
    return new ProteinWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;

