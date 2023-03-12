load csv with headers from "file:///uniqueMicroservices.csv" as row
create (n:MicroService{id:row.microservice})


load csv with headers from 'file:///graph.csv' as row
match (um:MicroService{id:row.um}), (dm:MicroService{id: row.dm})
create ((um)-[:RPC {timestamps: row.edgeValue}]->(dm))

