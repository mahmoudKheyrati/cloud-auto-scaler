#! /bin/bash

psql -h localhost -p 54320 -d auto-scaler -U admin -c "create table ms_call_graph
(
    id        bigint,
    trace_id  text,
    ts        text,
    rpc_id    text,
    um        text,
    rpc_type  text,
    dm        text,
    interface text,
    rt        int
);"
