#!/bin/zsh

# Run main.go in the background and capture its process name
go run main.go &
sleep 1
PROC_NAME=$(ps -e -o comm= | grep main)
echo "Process name: $PROC_NAME"

# Create a file to save the output
OUTPUT_FILE="performanceTEST.txt"
touch $OUTPUT_FILE

# Monitor the CPU and memory usage of the process with the captured PID
while true; do
  # Check if process name is empty
  if [ -z "$PROC_NAME" ]; then
    echo "Could not find process name."
    break
  fi
  
  # Get the PID of the process and check if it is empty
  PID=$(pgrep -f "exe/main")
  if [ -z "$PID" ]; then
    echo "Process $PROC_NAME is no longer running."
    break
  else
    ps -p $PID -o %cpu,%mem,command >> $OUTPUT_FILE
    sleep 0.5;
  fi
done
