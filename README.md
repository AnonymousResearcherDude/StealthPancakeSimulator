# StealthPancakeSimulator
This repository contains the primary version of the code referenced in the paper submitted to NSDI'24. Due to the double-blind review process, the repository is anonymized to maintain confidentiality.

## Instructions to compile and run
Ensure Golang, preferably version 1.19.5 or later, is installed on the computer.

After cloning the repository, configure the settings for the simulation by editing the `config.yaml` file in the root directory.

Before running the simulation, you need to generate a network in `generate_network_data` by running `go run generate_data.go`

Compile the code in the root directory with `go build`, that generates a binary file `StealthPancakeSimulator` and run it `./StealthPancakeSimulator`. 

You can find results of your simulation in the `results` directory. 
