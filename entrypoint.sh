#!/bin/sh

/jobstatus

sleep infinity & PID=$!
trap "kill $PID" INT TERM

echo sleeping

wait

echo exited
