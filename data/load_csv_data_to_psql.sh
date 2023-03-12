#! /bin/bash
psql -h localhost -p 54320 -d auto-scaler -U admin -c "\copy ms_call_graph from ./MSCallGraph_0.csv with delimiter as ',' csv header"

